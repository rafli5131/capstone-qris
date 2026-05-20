import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: Number(__ENV.VUS || 20),
  duration: __ENV.DURATION || '2m',
  thresholds: {
    http_req_duration: ['p(95)<1000'],
    http_req_failed: ['rate<0.50'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://legacy-system-java:8081';

http.setResponseCallback(http.expectedStatuses({ min: 200, max: 499 }));

export default function () {
  const res = http.get(`${BASE_URL}/api/admin/transactions`, {
    tags: { target_service: 'legacy-system-java', endpoint: 'admin_transactions' },
    headers: { Accept: 'application/json' },
  });

  check(res, {
    'legacy responded': (r) => r.status > 0,
    'legacy expected protected endpoint': (r) => [401, 403].includes(r.status),
  });

  sleep(1);
}
