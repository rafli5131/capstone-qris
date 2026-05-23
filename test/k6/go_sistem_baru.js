import http from 'k6/http';
import { check, sleep } from 'k6';
import { authHeaders, registerUser } from './auth_helpers.js';

export const options = {
  vus: Number(__ENV.VUS || 20),
  duration: __ENV.DURATION || '2m',
  thresholds: {
    http_req_duration: ['p(95)<1000'],
    http_req_failed: ['rate<0.01'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://go-sistem-baru:3000';

http.setResponseCallback(http.expectedStatuses({ min: 200, max: 399 }));

export function setup() {
  return {
    token: registerUser(BASE_URL, 'k6go'),
  };
}

export default function (data) {
  const res = http.get(`${BASE_URL}/api/admin/transactions`, {
    tags: { target_service: 'go-sistem-baru', endpoint: 'admin_transactions' },
    headers: authHeaders(data.token),
  });

  check(res, {
    'go responded': (r) => r.status > 0,
    'go status 200': (r) => r.status === 200,
  });

  sleep(1);
}
