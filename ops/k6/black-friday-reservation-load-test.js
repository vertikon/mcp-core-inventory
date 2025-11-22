import http from 'k6/http';
import { check, sleep, fail } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

// Custom metrics for Black Friday scenario
export let options = {
  stages: [
    { duration: '30s', target: 50 },   // Warm up - gradual ramp
    { duration: '1m', target: 200 },    // Black Friday start - medium traffic
    { duration: '2m', target: 500 },    // Peak Black Friday - high traffic
    { duration: '2m', target: 1000 },   // Extreme peak - flash sale moment
    { duration: '2m', target: 500 },    // Sustained high traffic
    { duration: '1m', target: 200 },    // Cool down
    { duration: '30s', target: 0 },     // Cool down complete
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],        // 95% of requests under 500ms
    http_req_failed: ['rate<0.1'],         // Less than 10% failure rate
    race_condition_incidents: ['count<5'],  // Maximum 5 race conditions
    stock_inconsistency: ['count<1'],      // Zero stock inconsistencies
  },
};

// Custom metrics
let raceConditionIncidents = new Counter('race_condition_incidents');
let stockInconsistency = new Counter('stock_inconsistency');
let reservationSuccess = new Counter('reservation_success');
let reservationFailure = new Counter('reservation_failure');
let insufficientStockErrors = new Counter('insufficient_stock_errors');
let reservationErrors = new Counter('reservation_errors');
let responseTimeTrend = new Trend('reservation_response_time');
let concurrentReservations = new Rate('concurrent_reservations');

// Configuration
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080/v1';
const BLACK_FRIDAY_SKU = __ENV.BLACK_FRIDAY_SKU || 'BLACK-FRIDAY-SPECIAL-001';
const WAREHOUSE_LOCATION = __ENV.WAREHOUSE_LOCATION || 'WH-MAIN';
const INITIAL_STOCK = parseInt(__ENV.INITIAL_STOCK) || 100;

// Test data pool for realistic load
const TEST_SKUS = [
  'BLACK-FRIDAY-SPECIAL-001', // High-demand item
  'SMARTPHONE-PRO-128GB',
  'LAPTOP-GAMING-RTX4080',
  'HEADPHONES-NOISE-CANCEL',
  'SMARTWATCH-FITNESS-PRO',
  'TABLET-11INCH-WIFI',
  'CAMERA-DSLR-4K',
  'SPEAKER-BLUETOOTH-PORTABLE',
  'MONITOR-27INCH-4K',
  'KEYBOARD-MECHANICAL-RGB'
];

// State tracking for race condition detection
let stockState = new Map();
let activeReservations = new Set();

export function setup() {
  console.log(`🎯 Black Friday Load Test Setup`);
  console.log(`📊 Target: ${BASE_URL}`);
  console.log(`🎁 Focus SKU: ${BLACK_FRIDAY_SKU}`);
  console.log(`📦 Expected Stock: ${INITIAL_STOCK}`);
  
  // Initialize stock tracking
  TEST_SKUS.forEach(sku => {
    stockState.set(sku, INITIAL_STOCK);
  });
  
  // Pre-warm the stock query endpoint
  let warmupResponse = http.get(`${BASE_URL}/available?sku=${BLACK_FRIDAY_SKU}&location=${WAREHOUSE_LOCATION}`);
  
  if (warmupResponse.status !== 200) {
    fail(`❌ Setup failed: Cannot query stock for ${BLACK_FRIDAY_SKU}. Status: ${warmupResponse.status}`);
  }
  
  let stockData = JSON.parse(warmupResponse.body);
  console.log(`✅ Initial stock for ${BLACK_FRIDAY_SKU}: ${stockData.stock.available}`);
  
  return { 
    initialStock: stockData.stock.available,
    testStartTime: Date.now()
  };
}

export default function(data) {
  // Simulate realistic user behavior: 70% try main item, 30% explore other items
  const targetSKU = Math.random() < 0.7 ? BLACK_FRIDAY_SKU : TEST_SKUS[randomIntBetween(1, TEST_SKUS.length - 1)];
  const quantity = Math.random() < 0.8 ? 1 : randomIntBetween(2, 3); // Most users buy 1, some buy multiple
  
  // Generate unique idempotency key
  const idempotencyKey = `black_friday_${Date.now()}_${__VU}_${randomIntBetween(1000, 9999)}`;
  
  // Start concurrent reservation attempt tracking
  const reservationStart = Date.now();
  activeReservations.add(idempotencyKey);
  
  let payload = JSON.stringify({
    sku: targetSKU,
    location: WAREHOUSE_LOCATION,
    quantity: quantity,
    idempotency_key: idempotencyKey,
    ttl: '15m' // 15 minutes TTL for realistic behavior
  });
  
  let params = {
    headers: {
      'Content-Type': 'application/json',
      'X-Test-Scenario': 'black-friday-load',
      'X-User-ID': `user_${__VU}`,
      'X-Request-ID': idempotencyKey
    },
  };
  
  // Execute reservation request
  let response = http.post(`${BASE_URL}/reserve`, payload, params);
  
  // Track response time
  const responseTime = Date.now() - reservationStart;
  responseTimeTrend.add(responseTime);
  activeReservations.delete(idempotencyKey);
  
  // Comprehensive response validation
  let success = check(response, {
    'reservation_status_is_200_or_409': (r) => r.status === 200 || r.status === 409,
    'reservation_response_time_under_1000ms': (r) => r.timings.duration < 1000,
    'reservation_has_proper_headers': (r) => r.headers['Content-Type'] === 'application/json',
    'reservation_json_valid': (r) => {
      try {
        if (r.body) JSON.parse(r.body);
        return true;
      } catch (e) {
        return false;
      }
    }
  });
  
  if (success) {
    let responseBody = JSON.parse(response.body);
    
    if (response.status === 200) {
      // Successful reservation - update stock state
      reservationSuccess.add(1);
      
      if (responseBody.reservation_id) {
        console.log(`✅ VU${__VU}: Reserved ${quantity} of ${targetSKU} (ID: ${responseBody.reservation_id})`);
        
        // Update local stock tracking for consistency checks
        let currentStock = stockState.get(targetSKU) || 0;
        stockState.set(targetSKU, currentStock - quantity);
      }
      
    } else if (response.status === 409 || response.status === 422) {
      // Expected business logic failure (insufficient stock)
      insufficientStockErrors.add(1);
      
      if (responseBody.error && responseBody.error.includes('insufficient')) {
        console.log(`⚠️  VU${__VU}: Insufficient stock for ${targetSKU}`);
      }
      
    } else {
      // Unexpected failure
      reservationFailure.add(1);
      console.log(`❌ VU${__VU}: Unexpected error ${response.status}: ${responseBody.error || 'Unknown error'}`);
    }
  } else {
    reservationErrors.add(1);
    console.log(`❌ VU${__VU}: Request failed with status ${response.status}`);
  }
  
  // Race condition detection - verify stock consistency
  if (Math.random() < 0.1) { // Check 10% of requests to avoid excessive load
    verifyStockConsistency(targetSKU);
  }
  
  // Simulate user think time - realistic browsing behavior
  sleep(randomIntBetween(100, 500) / 1000); // 100ms to 500ms
}

function verifyStockConsistency(sku) {
  let stockResponse = http.get(`${BASE_URL}/available?sku=${sku}&location=${WAREHOUSE_LOCATION}`);
  
  if (stockResponse.status === 200) {
    let stockData = JSON.parse(stockResponse.body);
    let trackedStock = stockState.get(sku) || 0;
    let actualStock = stockData.stock.available;
    
    // Check for stock inconsistencies (race conditions detected)
    if (actualStock < 0) {
      stockInconsistency.add(1);
      raceConditionIncidents.add(1);
      console.log(`🚨 RACE CONDITION: Negative stock detected for ${sku}: ${actualStock}`);
    }
    
    // Check for divergence between tracked and actual stock
    let divergence = Math.abs(trackedStock - actualStock);
    if (divergence > INITIAL_STOCK * 0.1) { // 10% tolerance
      raceConditionIncidents.add(1);
      console.log(`⚠️  STOCK DIVERGENCE for ${sku}: tracked=${trackedStock}, actual=${actualStock}`);
    }
  }
}

export function teardown(data) {
  const testDuration = (Date.now() - data.testStartTime) / 1000;
  
  console.log(`\n📊 Black Friday Load Test Results`);
  console.log(`⏱️  Test Duration: ${testDuration.toFixed(2)}s`);
  console.log(`✅ Successful Reservations: ${reservationSuccess.count}`);
  console.log(`❌ Failed Reservations: ${reservationFailure.count}`);
  console.log(`⚠️  Insufficient Stock Errors: ${insufficientStockErrors.count}`);
  console.log(`🔥 Race Condition Incidents: ${raceConditionIncidents.count}`);
  console.log(`🚨 Stock Inconsistencies: ${stockInconsistency.count}`);
  console.log(`📈 Average Response Time: ${responseTimeTrend.mean.toFixed(2)}ms`);
  console.log(`⚡ 95th Percentile: ${responseTimeTrend.p(95).toFixed(2)}ms`);
  
  // Black Friday success criteria
  const successRate = reservationSuccess.count / (reservationSuccess.count + reservationFailure.count) * 100;
  const raceConditionRate = raceConditionIncidents.count / testDuration * 60; // per minute
  
  console.log(`\n🎯 Black Friday Assessment`);
  console.log(`✅ Success Rate: ${successRate.toFixed(2)}%`);
  console.log(`🔥 Race Conditions: ${raceConditionRate.toFixed(2)}/min`);
  
  if (raceConditionIncidents.count === 0 && stockInconsistency.count === 0) {
    console.log(`🏆 EXCELLENT: No race conditions detected! System handled Black Friday perfectly.`);
  } else if (raceConditionIncidents.count < 5) {
    console.log(`✅ GOOD: Minimal race conditions. System performed well under extreme load.`);
  } else {
    console.log(`⚠️  ATTENTION: ${raceConditionIncidents.count} race conditions detected. Review locking mechanism.`);
  }
  
  // Final stock verification
  let finalStockResponse = http.get(`${BASE_URL}/available?sku=${BLACK_FRIDAY_SKU}&location=${WAREHOUSE_LOCATION}`);
  if (finalStockResponse.status === 200) {
    let finalStock = JSON.parse(finalStockResponse.body);
    console.log(`\n📦 Final Stock Status for ${BLACK_FRIDAY_SKU}:`);
    console.log(`   Available: ${finalStock.stock.available}`);
    console.log(`   Reserved: ${finalStock.stock.reserved}`);
    console.log(`   Total: ${finalStock.stock.total}`);
    
    if (finalStock.stock.available < 0) {
      console.log(`🚨 CRITICAL: Negative stock detected! This is a race condition bug.`);
    }
  }
}

// Utility function for generating realistic user behavior patterns
export function handleReservationsFlow() {
  // Simulate typical Black Friday user journey:
  // 1. Check stock availability
  // 2. Attempt reservation
  // 3. If successful, proceed to confirm (5% of cases)
  // 4. If failed, try alternative product (20% of cases)
  
  // This function can be called from the main test function for more complex scenarios
}

export function simulateFlashSale() {
  // Simulate flash sale scenario - sudden spike in traffic
  // This would be used in a separate test stage
  console.log(`⚡ FLASH SALE MODE: Simulating sudden traffic spike`);
  
  // Target the high-demand item specifically
  let flashSalePayload = {
    sku: BLACK_FRIDAY_SKU,
    location: WAREHOUSE_LOCATION,
    quantity: 1,
    idempotency_key: `flash_sale_${Date.now()}_${__VU}`,
    ttl: '5m' // Shorter TTL for flash sales
  };
  
  return JSON.stringify(flashSalePayload);
}