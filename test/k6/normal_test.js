import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';

const goHttpReqs = new Counter('go_http_reqs');
const goHttpReqFailed = new Rate('go_http_req_failed');
const goHttpReqDuration = new Trend('go_http_req_duration', true);
const legacyHttpReqs = new Counter('legacy_http_reqs');
const legacyHttpReqFailed = new Rate('legacy_http_req_failed');
const legacyHttpReqDuration = new Trend('legacy_http_req_duration', true);

export const options = {
  scenarios: {
    go_sistem_baru_normal: {
      executor: 'ramping-vus',
      exec: 'goSistemBaru',
      stages: [
        { duration: '30s', target: Number(__ENV.GO_VUS || 100) },
        { duration: '120s', target: Number(__ENV.GO_VUS || 100) },
        { duration: '30s', target: 0 },
      ],
      tags: { target_service: 'go-sistem-baru', test_scenario: 'normal' },
    },
    legacy_system_java_normal: {
      executor: 'ramping-vus',
      exec: 'legacySystemJava',
      stages: [
        { duration: '30s', target: Number(__ENV.LEGACY_VUS || 100) },
        { duration: '120s', target: Number(__ENV.LEGACY_VUS || 100) },
        { duration: '30s', target: 0 },
      ],
      tags: { target_service: 'legacy-system-java', test_scenario: 'normal' },
    },
  },
  thresholds: {
    'http_req_duration{test_scenario:normal}': ['p(95)<500'],
    'http_req_failed{test_scenario:normal}': ['rate<0.01'],
    go_http_req_duration: ['p(95)<500'],
    go_http_req_failed: ['rate<0.01'],
    legacy_http_req_duration: ['p(95)<500'],
    legacy_http_req_failed: ['rate<0.01'],
  },
};

const GO_BASE_URL = __ENV.GO_BASE_URL || 'http://go-sistem-baru:3000';
const LEGACY_BASE_URL = __ENV.LEGACY_BASE_URL || 'http://legacy-system-java:8081';

http.setResponseCallback(http.expectedStatuses({ min: 200, max: 499 }));

function exerciseProtectedEndpoint(baseUrl, serviceName) {
  group(`${serviceName} - Normal Network`, () => {
    const res = http.get(`${baseUrl}/api/admin/transactions`, {
      tags: { endpoint: 'admin_transactions' },
      headers: { Accept: 'application/json' },
    });

    check(res, {
      [`${serviceName} responded`]: (r) => r.status > 0,
      [`${serviceName} protected endpoint`]: (r) => [401, 403].includes(r.status),
      [`${serviceName} response < 500ms`]: (r) => r.timings.duration < 500,
    });

    recordServiceMetrics(serviceName, res);
  });

  sleep(0.5);
}

function recordServiceMetrics(serviceName, res) {
  const failed = res.status >= 500 || res.status === 0;
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

export function goSistemBaru() {
  exerciseProtectedEndpoint(GO_BASE_URL, 'go-sistem-baru');
}

export function legacySystemJava() {
  exerciseProtectedEndpoint(LEGACY_BASE_URL, 'legacy-system-java');
}
