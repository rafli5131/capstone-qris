# Load Testing Documentation

This document describes the load testing scenarios for the QRIS Payment System API.

## Test Scenarios

### 1. Rural Network Test (`rural_test.js`)
- **Simulates**: Low bandwidth, high latency rural network conditions
- **Concurrent Users**: 10
- **Latency Simulation**: 0.5-2s random delays between requests
- **Duration**: 180 seconds (30s ramp-up, 120s steady, 30s ramp-down)
- **Thresholds**:
  - 95% of requests < 3s
  - Error rate < 5%

### 2. Normal Network Test (`normal_test.js`)
- **Simulates**: Standard urban network conditions
- **Concurrent Users**: 100
- **Latency Simulation**: Minimal delays (0.5s between requests)
- **Duration**: 180 seconds
- **Thresholds**:
  - 95% of requests < 500ms
  - Error rate < 1%

### 3. Peak Load Test (`peak_test.js`)
- **Simulates**: High traffic peak hours
- **Concurrent Users**: 1000
- **Latency Simulation**: Short delays (0.2-0.5s between requests)
- **Duration**: 180 seconds
- **Thresholds**:
  - 95% of requests < 1s
  - Error rate < 5%

### 4. Peak Rural Network Test (`peak_rural_test.js`)
- **Simulates**: Peak load under rural network conditions
- **Concurrent Users**: 500
- **Latency Simulation**: 0.5-1.5s random delays
- **Duration**: 180 seconds
- **Thresholds**:
  - 95% of requests < 5s
  - Error rate < 10%

## Running Tests

### Prerequisites
- k6 installed (`brew install k6` or download from k6.io)
- API server running on localhost:3000 (or set BASE_URL env var)
- Valid JWT token (set JWT_TOKEN env var or use default mock token)

### Commands

```bash
# Rural Network Test
k6 run test/k6/rural_test.js

# Normal Network Test
k6 run test/k6/normal_test.js

# Peak Load Test
k6 run test/k6/peak_test.js

# Peak Rural Network Test
k6 run test/k6/peak_rural_test.js

# With custom base URL and JWT token
BASE_URL=http://your-api-server:3000 JWT_TOKEN=your-jwt-token-here k6 run test/k6/normal_test.js
```

## Test Flow

Each test scenario performs these operations in sequence:

1. **QRIS Inquiry**: GET request to `/api/qris/inquiry/{qris_payload}`
   - Uses JWT token in Authorization header
   - Validates 200 status
   - Checks for metadata in response
   - Measures response time

2. **QRIS Payment**: POST request to `/api/qris/payment`
   - Uses JWT token in Authorization header
   - Sends payment request with random amount
   - Validates 202 status
   - Checks for transaction_id in response
   - Uses same JWT token as inquiry

## Monitoring

During test execution, k6 will output:
- Real-time metrics (requests/sec, response times)
- Progress indicators
- Final summary with pass/fail status

## Expected Results

- **Rural**: Higher response times acceptable, focus on stability
- **Normal**: Fast responses, high reliability
- **Peak**: Good performance under load, acceptable error rate
- **Peak Rural**: Balanced performance considering network constraints

## Troubleshooting

- If tests fail due to authentication, ensure JWT_TOKEN is valid and not expired
- For connection errors, verify API server is running and accessible
- High error rates may indicate server capacity issues or network problems
- Default mock JWT token is provided but may be expired - replace with a fresh token from login endpoint