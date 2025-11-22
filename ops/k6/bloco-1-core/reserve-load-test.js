import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

// --- CONFIGURAÇÕES DE TESTE BLOCO-1 ---
export const options = {
  scenarios: {
    black_friday_bloco1: {
      executor: 'ramping-arrival-rate',
      startRate: 0,
      timeUnit: '1s',
      preAllocatedVUs: 50,
      maxVUs: 1000, // Simulando alta concorrência
      stages: [
        { target: 50, duration: '30s' },  // Aquecimento
        { target: 200, duration: '1m' },  // Pico de vendas (Flash Sale)
        { target: 500, duration: '2m' },  // Pico extremo (Black Friday)
        { target: 1000, duration: '3m' }, // Estresse máximo
        { target: 0, duration: '30s' },   // Esfriamento
      ],
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.01'], // Erros 500 devem ser quase zero
    http_req_duration: ['p95<200'], // 95% das requisições abaixo de 200ms
    'business_errors': ['count>0'], // Queremos ver erros 409/422 (significa que o bloqueio funcionou!)
    'race_conditions': ['count=0'], // ZERO race conditions é o objetivo
  },
};

// --- MÉTRICAS CUSTOMIZADAS BLOCO-1 ---
const businessErrors = new Counter('business_errors'); // Conta conflitos (estoque acabou)
const successfulReservations = new Counter('successful_reservations');
const raceConditionSuspects = new Counter('race_condition_suspects');
const inventoryLatency = new Trend('inventory_reservation_latency', true);
const concurrentUsers = new Trend('concurrent_users', true);

// --- SETUP (Opcional: Criar massa de dados antes) ---
export function setup() {
    // Exemplo: Garantir saldo 100 para o SKU "VESTIDO-VERAO-P"
    // http.post(...) 
    return { 
        sku: "VESTIDO-VERAO-P", 
        location: "CD-SP",
        initial_stock: 100
    };
}

// --- TESTE PRINCIPAL BLOCO-1 ---
export default function (data) {
  const baseUrl = __ENV.BASE_URL || 'http://localhost:8080';
  const url = `${baseUrl}/v1/reserve`;
  
  // Payload simulando uma compra no BLOCO-1
  const payload = JSON.stringify({
    sku: data.sku,
    location: data.location,
    quantity: 1,
    // Importante: Idempotency Key única por request para forçar processamento real
    idempotency_key: uuidv4(), 
    ttl_seconds: 300,
    // Metadados do BLOCO-1
    metadata: {
      source: 'k6-stress-test',
      bloco: '1',
      test_type: 'black_friday'
    }
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'X-Tenant-ID': 'vertikon-demo', // Se usar multitenancy
      'X-Bloco-ID': '1', // Identificador do BLOCO-1
      'X-Test-Type': 'stress-test'
    },
    tags: { 
        name: 'ReserveStock',
        bloco: '1',
        service: 'core-inventory'
    },
  };

  const startTime = new Date().getTime();
  const res = http.post(url, payload, params);
  const endTime = new Date().getTime();
  const responseTime = endTime - startTime;

  // Registrar latência
  inventoryLatency.add(responseTime);

  // --- VALIDAÇÕES BLOCO-1 ---
  
  // 1. Sucesso Técnico (O servidor respondeu?)
  check(res, {
    'status is valid (200, 409, 422)': (r) => [200, 201, 409, 422].includes(r.status),
    'response time < 200ms': (r) => responseTime < 200,
    'content-type is application/json': (r) => r.headers['Content-Type'] && r.headers['Content-Type'].includes('application/json'),
  });

  // 2. Lógica de Negócio (O que aconteceu com o estoque?)
  if (res.status === 200 || res.status === 201) {
    successfulReservations.add(1);
  } else if (res.status === 409 || res.status === 422) {
    // 409 = Conflict (Lock falhou ou concorrencia)
    // 422 = Unprocessable (Sem saldo)
    businessErrors.add(1); 
  } else {
    // 500 = Erro de código/banco (ISSO NÃO PODE ACONTECER)
    raceConditionSuspects.add(1);
    console.error(`Erro Crítico BLOCO-1 ${res.status}: ${res.body}`);
  }

  sleep(0.1); // Pequena pausa para não DDOSar a rede local
}

// --- TEARDOWN BLOCO-1 ---
export function teardown(data) {
    console.log('=== BLOCO-1 CORE INVENTORY TEST RESULTS ===');
    console.log(`SKU Testado: ${data.sku}`);
    console.log(`Localização: ${data.location}`);
    console.log(`Saldo Inicial: ${data.initial_stock}`);
    console.log(`Reservas Bem-sucedidas: ${successfulReservations.count}`);
    console.log(`Erros de Negócio: ${businessErrors.count}`);
    console.log(`Suspeitas de Race Condition: ${raceConditionSuspects.count}`);
    console.log(`Latência Média: ${inventoryLatency.mean}ms`);
    console.log(`Latência p95: ${inventoryLatency.p(95)}ms`);
    console.log(`Latência p99: ${inventoryLatency.p(99)}ms`);
    
    // Validação crítica do BLOCO-1
    if (successfulReservations.count > data.initial_stock) {
        console.error('🚨 CRÍTICO: OVERSELLING DETECTADO! Vendas > Saldo Inicial');
    } else {
        console.log('✅ OK: Zero Overselling - BLOCO-1 funcionando corretamente');
    }
    
    if (raceConditionSuspects.count > 0) {
        console.error('🚨 CRÍTICO: RACE CONDITIONS DETECTADAS!');
    } else {
        console.log('✅ OK: Zero Race Conditions - BLOCO-1 estável sob concorrência');
    }
}
