import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { authHeaders, registerUser } from './auth_helpers.js';

export const options = {
  scenarios: {
    go_sistem_baru: {
      executor: 'constant-vus',
      vus: Number(__ENV.GO_VUS || 20),
      duration: __ENV.DURATION || '2m',
      exec: 'goSistemBaru',
      tags: { target_service: 'go-sistem-baru' },
    },
    legacy_system_java: {
      executor: 'constant-vus',
      vus: Number(__ENV.LEGACY_VUS || 20),
      duration: __ENV.DURATION || '2m',
      exec: 'legacySystemJava',
      tags: { target_service: 'legacy-system-java' },
    },
  },
  thresholds: {
    'http_req_duration{target_service:go-sistem-baru}': ['p(95)<1000'],
    'http_req_duration{target_service:legacy-system-java}': ['p(95)<1000'],
    http_req_failed: ['rate<0.01'],
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

export function goSistemBaru(data) {
  group('go_sistem_baru', () => {
    const res = http.get(`${GO_BASE_URL}/api/admin/transactions`, {
      tags: { app: 'go-sistem-baru', endpoint: 'admin_transactions' },
      headers: authHeaders(data.goToken),
    });

    check(res, {
      'go responded': (r) => r.status > 0,
      'go status 200': (r) => r.status === 200,
    });
  });

  sleep(1);
}

export function legacySystemJava(data) {
  group('legacy_system_java', () => {
    const res = http.get(`${LEGACY_BASE_URL}/api/admin/transactions`, {
      tags: { app: 'legacy-system-java', endpoint: 'admin_transactions' },
      headers: authHeaders(data.legacyToken),
    });

    check(res, {
      'legacy responded': (r) => r.status > 0,
      'legacy status 200': (r) => r.status === 200,
    });
  });

  sleep(1);
}
