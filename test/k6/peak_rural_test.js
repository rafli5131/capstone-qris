import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';
import { authHeaders, registerUser } from './auth_helpers.js';

const goHttpReqs = new Counter('go_http_reqs');
const goHttpReqFailed = new Rate('go_http_req_failed');
const goHttpReqDuration = new Trend('go_http_req_duration', true);
const legacyHttpReqs = new Counter('legacy_http_reqs');
const legacyHttpReqFailed = new Rate('legacy_http_req_failed');
const legacyHttpReqDuration = new Trend('legacy_http_req_duration', true);

export const options = {
  scenarios: {
    go_sistem_baru_peak_rural: {
      executor: 'ramping-vus',
      exec: 'goSistemBaru',
      stages: [
        { duration: '30s', target: Number(__ENV.GO_VUS || 500) },
        { duration: '120s', target: Number(__ENV.GO_VUS || 500) },
        { duration: '30s', target: 0 },
      ],
      tags: { target_service: 'go-sistem-baru', test_scenario: 'peak_rural' },
    },
    legacy_system_java_peak_rural: {
      executor: 'ramping-vus',
      exec: 'legacySystemJava',
      stages: [
        { duration: '30s', target: Number(__ENV.LEGACY_VUS || 500) },
        { duration: '120s', target: Number(__ENV.LEGACY_VUS || 500) },
        { duration: '30s', target: 0 },
      ],
      tags: { target_service: 'legacy-system-java', test_scenario: 'peak_rural' },
    },
  },
  thresholds: {
    'http_req_duration{test_scenario:peak_rural}': ['p(95)<5000'],
    'http_req_failed{test_scenario:peak_rural}': ['rate<0.10'],
    go_http_req_duration: ['p(95)<5000'],
    go_http_req_failed: ['rate<0.10'],
    legacy_http_req_duration: ['p(95)<5000'],
    legacy_http_req_failed: ['rate<0.10'],
  },
};

const GO_BASE_URL = __ENV.GO_BASE_URL || 'http://go-sistem-baru:3000';
const LEGACY_BASE_URL = __ENV.LEGACY_BASE_URL || 'http://legacy-system-java:8081';

http.setResponseCallback(http.expectedStatuses({ min: 200, max: 399 }));

export function setup() {
  return {
    goToken: registerUser(GO_BASE_URL, 'k6go'),
    legacyToken: registerUser(LEGACY_BASE_URL, 'k6legacy'),
  };
}

function exerciseProtectedEndpoint(baseUrl, serviceName, token) {
  sleep(Math.random() + 0.5);

  group(`${serviceName} - Peak Rural Network`, () => {
    const res = http.get(`${baseUrl}/api/admin/transactions`, {
      tags: { endpoint: 'admin_transactions' },
      headers: authHeaders(token),
    });

    check(res, {
      [`${serviceName} responded`]: (r) => r.status > 0,
      [`${serviceName} status 200`]: (r) => r.status === 200,
      [`${serviceName} response < 5s`]: (r) => r.timings.duration < 5000,
    });

    recordServiceMetrics(serviceName, res);
  });
}

function recordServiceMetrics(serviceName, res) {
  const failed = res.status !== 200;
  if (serviceName === 'go-sistem-baru') {
    goHttpReqs.add(1);
    goHttpReqFailed.add(failed);
    goHttpReqDuration.add(res.timings.duration);
    return;
  }

  legacyHttpReqs.add(1);
  legacyHttpReqFailed.add(failed);
  legacyHttpReqDuration.add(res.timings.duration);
}

export function goSistemBaru(data) {
  exerciseProtectedEndpoint(GO_BASE_URL, 'go-sistem-baru', data.goToken);
}

export function legacySystemJava(data) {
  exerciseProtectedEndpoint(LEGACY_BASE_URL, 'legacy-system-java', data.legacyToken);
}
