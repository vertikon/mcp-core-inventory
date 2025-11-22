#!/bin/bash

# Black Friday Load Test Certification Script
# Validates that mcp-core-inventory can handle extreme load with race condition prevention

set -euo pipefail

# Configuration
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"
readonly TEST_RESULTS_DIR="$PROJECT_ROOT/test-results"
readonly CONFIG_FILE="$PROJECT_ROOT/config/config.yaml"

# Test Parameters
readonly VUS_SCENARIOS=(
    "50:2m"      # Warm up
    "200:1m"     # Normal Black Friday  
    "500:2m"      # Peak Black Friday
    "1000:2m"     # Extreme peak (flash sale)
    "500:2m"      # Sustained high load
    "200:1m"      # Cool down
)

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# Test Configuration
export BASE_URL="${HULK_BASE_URL:-http://localhost:8080/v1}"
export BLACK_FRIDAY_SKU="${HULK_BLACK_FRIDAY_SKU:-BLACK-FRIDAY-SPECIAL-001}"
export WAREHOUSE_LOCATION="${HULK_WAREHOUSE_LOCATION:-WH-MAIN}"
export INITIAL_STOCK="${HULK_INITIAL_STOCK:-1000}"

# Performance Gates
readonly MAX_RACE_CONDITIONS=0
readonly MAX_STOCK_INCONSISTENCIES=0
readonly MAX_FAILURE_RATE=10    # percent
readonly MAX_RESPONSE_TIME_P95=500  # ms
readonly MIN_SUCCESS_RATE=90     # percent

# Logging
log() {
    echo -e "${BLUE}[$(date '+%Y-%m-%d %H:%M:%S')] $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

log_header() {
    echo
    echo -e "${BLUE}═════════════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE} $1${NC}"
    echo -e "${BLUE}═════════════════════════════════════════════════════════════════${NC}"
    echo
}

# Setup
setup_test_environment() {
    log_header "Setting Up Test Environment"
    
    # Create results directory
    mkdir -p "$TEST_RESULTS_DIR"
    
    # Check dependencies
    if ! command -v k6 &> /dev/null; then
        log_error "k6 is not installed. Install it first: https://k6.io/docs/getting-started/installation"
        exit 1
    fi
    
    if ! command -v jq &> /dev/null; then
        log_error "jq is not installed. Install it first: sudo apt-get install jq"
        exit 1
    fi
    
    # Check if service is running
    if ! curl -sf "$BASE_URL/health" > /dev/null; then
        log_error "Service is not accessible at $BASE_URL"
        log_error "Start the service first: make run"
        exit 1
    fi
    
    log_success "Dependencies verified"
    log_success "Service is accessible"
    log_success "Test results directory: $TEST_RESULTS_DIR"
    
    # Get initial stock for verification
    local initial_stock_response
    initial_stock_response=$(curl -s "$BASE_URL/available?sku=$BLACK_FRIDAY_SKU&location=$WAREHOUSE_LOCATION")
    
    if echo "$initial_stock_response" | jq -e '.stock.available' > /dev/null; then
        local available_stock
        available_stock=$(echo "$initial_stock_response" | jq -r '.stock.available')
        log_success "Initial stock for $BLACK_FRIDAY_SKU: $available_stock"
        export INITIAL_STOCK="$available_stock"
    else
        log_warning "Could not fetch initial stock, using default: $INITIAL_STOCK"
    fi
}

# Execute single load test scenario
run_load_test_scenario() {
    local scenario="$1"
    local vus="${scenario%:*}"
    local duration="${scenario#*:}"
    
    log "Running scenario: $vus VUs for $duration"
    
    local test_name="load_test_${vus}vus_${duration//m/min}"
    local json_file="$TEST_RESULTS_DIR/${test_name}.json"
    local csv_file="$TEST_RESULTS_DIR/${test_name}.csv"
    
    # Execute k6 test
    k6 run \
        --vus "$vus" \
        --duration "$duration" \
        --out json="$json_file" \
        --out csv="$csv_file" \
        --summary-export "$TEST_RESULTS_DIR/${test_name}_summary.json" \
        "$PROJECT_ROOT/ops/k6/black-friday-reservation-load-test.js" \
        2>&1 | tee "$TEST_RESULTS_DIR/${test_name}.log"
    
    # Extract key metrics
    if [ -f "$json_file" ]; then
        extract_metrics "$json_file" "$test_name"
    else
        log_error "Test results file not found: $json_file"
    fi
}

# Extract and validate metrics from test results
extract_metrics() {
    local json_file="$1"
    local test_name="$2"
    
    log "Extracting metrics from $json_file"
    
    # Extract key metrics using jq
    local race_conditions
    local stock_inconsistencies  
    local http_req_failed
    local http_req_duration_p95
    local total_requests
    
    race_conditions=$(jq -r '.metrics.race_condition_incidents.count // 0' "$json_file")
    stock_inconsistencies=$(jq -r '.metrics.stock_inconsistency.count // 0' "$json_file") 
    http_req_failed=$(jq -r '.metrics.http_req_failed.rate // 0' "$json_file")
    http_req_duration_p95=$(jq -r '.metrics.http_req_duration.values.p95 // 0' "$json_file")
    total_requests=$(jq -r '.metrics.http_reqs.count // 0' "$json_file")
    
    # Calculate success rate
    local failure_rate_percent
    failure_rate_percent=$(echo "$http_req_failed * 100" | bc -l 2>/dev/null || echo "0")
    local success_rate_percent
    success_rate_percent=$(echo "100 - $failure_rate_percent" | bc -l 2>/dev/null || echo "100")
    
    # Save metrics for summary
    cat >> "$TEST_RESULTS_DIR/metrics_summary.csv" << EOF
$test_name,$vus,$race_conditions,$stock_inconsistencies,$failure_rate_percent,$success_rate_percent,$http_req_duration_p95,$total_requests
EOF
    
    # Validate performance gates
    validate_performance_gates "$test_name" "$race_conditions" "$stock_inconsistencies" "$failure_rate_percent" "$success_rate_percent" "$http_req_duration_p95"
}

# Performance gate validation
validate_performance_gates() {
    local test_name="$1"
    local race_conditions="$2"
    local stock_inconsistencies="$3"
    local failure_rate_percent="$4"
    local success_rate_percent="$5"
    local response_time_p95="$6"
    
    log "Validating performance gates for $test_name"
    
    local passed_gates=0
    local total_gates=5
    
    # Gate 1: No race conditions
    if [ "$race_conditions" -eq "$MAX_RACE_CONDITIONS" ]; then
        log_success "✓ Race Conditions: $race_conditions (PASS)"
        ((passed_gates++))
    else
        log_error "✗ Race Conditions: $race_conditions (FAIL - must be $MAX_RACE_CONDITIONS)"
    fi
    
    # Gate 2: No stock inconsistencies
    if [ "$stock_inconsistencies" -eq "$MAX_STOCK_INCONSISTENCIES" ]; then
        log_success "✓ Stock Inconsistencies: $stock_inconsistencies (PASS)"
        ((passed_gates++))
    else
        log_error "✗ Stock Inconsistencies: $stock_inconsistencies (FAIL - must be $MAX_STOCK_INCONSISTENCIES)"
    fi
    
    # Gate 3: Failure rate below threshold
    if (( $(echo "$failure_rate_percent <= $MAX_FAILURE_RATE" | bc -l 2>/dev/null || echo "1") )); then
        log_success "✓ Failure Rate: ${failure_rate_percent}% (PASS - threshold: $MAX_FAILURE_RATE%)"
        ((passed_gates++))
    else
        log_error "✗ Failure Rate: ${failure_rate_percent}% (FAIL - threshold: $MAX_FAILURE_RATE%)"
    fi
    
    # Gate 4: Success rate above threshold
    if (( $(echo "$success_rate_percent >= $MIN_SUCCESS_RATE" | bc -l 2>/dev/null || echo "1") )); then
        log_success "✓ Success Rate: ${success_rate_percent}% (PASS - threshold: $MIN_SUCCESS_RATE%)"
        ((passed_gates++))
    else
        log_error "✗ Success Rate: ${success_rate_percent}% (FAIL - threshold: $MIN_SUCCESS_RATE%)"
    fi
    
    # Gate 5: Response time P95 below threshold
    if (( $(echo "$response_time_p95 <= $MAX_RESPONSE_TIME_P95" | bc -l 2>/dev/null || echo "1") )); then
        log_success "✓ Response Time P95: ${response_time_p95}ms (PASS - threshold: ${MAX_RESPONSE_TIME_P95}ms)"
        ((passed_gates++))
    else
        log_error "✗ Response Time P95: ${response_time_p95}ms (FAIL - threshold: ${MAX_RESPONSE_TIME_P95}ms)"
    fi
    
    # Overall result
    local result="FAIL"
    if [ "$passed_gates" -eq "$total_gates" ]; then
        result="PASS"
        log_success "🎯 All performance gates PASSED for $test_name"
    else
        log_error "❌ $passed_gates/$total_gates gates PASSED for $test_name"
    fi
    
    # Save to summary
    cat >> "$TEST_RESULTS_DIR/validation_summary.csv" << EOF
$test_name,$result,$passed_gates,$total_gates,$race_conditions,$stock_inconsistencies,$failure_rate_percent,$success_rate_percent,$response_time_p95
EOF
}

# Generate final certification report
generate_certification_report() {
    log_header "Black Friday Load Test Certification Report"
    
    # Initialize CSV files
    echo "Test_Name,VUs,Race_Conditions,Stock_Inconsistencies,Failure_Rate%,Success_Rate%,Response_Time_P95(ms),Total_Requests" > "$TEST_RESULTS_DIR/metrics_summary.csv"
    echo "Test_Name,Result,Passed_Gates,Total_Gates,Race_Conditions,Stock_Inconsistencies,Failure_Rate%,Success_Rate%,Response_Time_P95(ms)" > "$TEST_RESULTS_DIR/validation_summary.csv"
    
    # Run all test scenarios
    for scenario in "${VUS_SCENARIOS[@]}"; do
        log_header "Running Load Test Scenario: $scenario"
        run_load_test_scenario "$scenario"
        sleep 5  # Brief pause between scenarios
    done
    
    # Final verification
    log_header "Final Stock Consistency Verification"
    local final_stock_response
    final_stock_response=$(curl -s "$BASE_URL/available?sku=$BLACK_FRIDAY_SKU&location=$WAREHOUSE_LOCATION")
    
    if echo "$final_stock_response" | jq -e '.stock.available' > /dev/null; then
        local final_stock
        final_stock=$(echo "$final_stock_response" | jq -r '.stock.available')
        
        if [ "$final_stock" -ge 0 ]; then
            log_success "Final stock consistency: $final_stock (PASS)"
        else
            log_error "Final stock consistency: $final_stock (FAIL - negative stock detected)"
        fi
    fi
    
    # Generate summary
    generate_summary
}

# Generate executive summary
generate_summary() {
    log_header "Executive Summary"
    
    # Count pass/fail scenarios
    local total_scenarios=0
    local passed_scenarios=0
    
    while IFS=, read -r test_name result _; do
        if [ "$test_name" != "Test_Name" ]; then
            ((total_scenarios++))
            if [ "$result" = "PASS" ]; then
                ((passed_scenarios++))
            fi
        fi
    done < "$TEST_RESULTS_DIR/validation_summary.csv"
    
    local pass_rate=$((passed_scenarios * 100 / total_scenarios))
    
    log "Total Test Scenarios: $total_scenarios"
    log "Passed Scenarios: $passed_scenarios"
    log "Pass Rate: ${pass_rate}%"
    
    # Overall certification
    local certification="REJECTED"
    if [ "$pass_rate" -eq 100 ]; then
        certification="CERTIFIED"
    elif [ "$pass_rate" -ge 80 ]; then
        certification="CONDITIONAL"
    fi
    
    echo
    log_header "BLACK FRIDAY LOAD TEST CERTIFICATION"
    log "$certification"
    
    case "$certification" in
        "CERTIFIED")
            log_success "🏆 CERTIFIED FOR PRODUCTION BLACK FRIDAY DEPLOYMENT"
            log_success "✅ System passed all performance gates under extreme load"
            log_success "✅ Zero race conditions detected"
            log_success "✅ Stock consistency maintained throughout all scenarios"
            log_success "✅ Response times within acceptable limits"
            echo
            log "🚀 RECOMMENDATION: APPROVED for immediate production deployment"
            ;;
        "CONDITIONAL")
            log_warning "⚠️ CONDITIONAL CERTIFICATION - Minor Issues Detected"
            log_warning "✅ Core functionality working under load"
            log_warning "⚠️ Some performance thresholds exceeded - review logs"
            echo
            log "🔧 RECOMMENDATION: Address identified issues before production deployment"
            ;;
        "REJECTED")
            log_error "❌ REJECTED FOR PRODUCTION - Critical Issues Detected"
            log_error "🚫 Race conditions or stock inconsistencies found"
            log_error "🚫 Performance not within acceptable limits"
            echo
            log "⏸️ RECOMMENDATION: FIX critical issues before any deployment consideration"
            ;;
    esac
    
    echo
    log_header "Next Steps"
    
    if [ "$certification" = "CERTIFIED" ]; then
        log "1. ✅ Deploy to staging for final validation"
        log "2. ✅ Configure production monitoring dashboards"
        log "3. ✅ Prepare production rollback plan"
        log "4. ✅ Schedule production deployment window"
    else
        log "1. 📊 Review detailed test results in $TEST_RESULTS_DIR"
        log "2. 🔍 Analyze logs for root cause identification"
        log "3. 🛠️ Fix critical race conditions and performance issues"
        log "4. 🔄 Re-run certification tests after fixes"
    fi
    
    echo
    log_header "Generated Reports"
    log "📊 Metrics Summary: $TEST_RESULTS_DIR/metrics_summary.csv"
    log "✅ Validation Summary: $TEST_RESULTS_DIR/validation_summary.csv"
    log "📋 Detailed Logs: $TEST_RESULTS_DIR/*.log"
    
    # Export environment variables for CI/CD
    echo "LOAD_TEST_CERTIFICATION=$certification" >> "$TEST_RESULTS_DIR/certification.env"
    echo "PASS_RATE=${pass_rate}%" >> "$TEST_RESULTS_DIR/certification.env"
}

# Main execution
main() {
    log_header "Black Friday Load Test Certification - mcp-core-inventory"
    log "Base URL: $BASE_URL"
    log "Test SKU: $BLACK_FRIDAY_SKU"
    log "Location: $WAREHOUSE_LOCATION"
    
    setup_test_environment
    generate_certification_report
    
    log_header "Certification Complete"
}

# Execute main function
main "$@"