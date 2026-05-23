import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { authHeaders, registerUser } from './auth_helpers.js';

export const options = {
  stages: [
    { duration: '30s', target: 1000 },
    { duration: '120s', target: 1000 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';
const SAMPLE_QRIS = '00020101021126230219ID102327560673003035204274153033605802ID5912M%20Ivan%20Store6013Jakarta%20Timur6304ABCD';

http.setResponseCallback(http.expectedStatuses({ min: 200, max: 399 }));

export function setup() {
  return {
    token: registerUser(BASE_URL, 'k6flow'),
  };
}

export default function (data) {
  let transactionIds = []; // Store transaction IDs from payments
  let inquiryId = null;

  group('QRIS Inquiry', () => {
    const res = http.get(`${BASE_URL}/api/qris/inquiry/${SAMPLE_QRIS}`, {
      headers: authHeaders(data.token),
    });
    check(res, {
      'status 200': (r) => r.status === 200,
      'has metadata': (r) => r.json('metadata') !== undefined,
      'has inquiry_id': (r) => r.json('data.inquiry_id') !== undefined,
    });
    if (res.status === 200) {
      inquiryId = res.json('data.inquiry_id');
    }
  });

  sleep(0.5);

  // Create 5 payments and store their transaction IDs
  for (let i = 0; i < 5; i++) {
    group(`QRIS Payment ${i + 1}`, () => {
      const body = JSON.stringify({
        inquiry_id: inquiryId,
        amount: Math.floor(Math.random() * 100000) + 1000,
        payment_method: 'balance',
        pincode: '123456',
      });
      const res = http.post(`${BASE_URL}/api/qris/payment`, body, {
        headers: {
          ...authHeaders(data.token),
          'Content-Type': 'application/json',
        },
      });
      check(res, {
        'status 202': (r) => r.status === 202,
        'has transaction_id': (r) => r.json('data.transaction_id') !== undefined,
      });

      // Store transaction ID for later status checks
      if (res.status === 202 && res.json('data.transaction_id')) {
        transactionIds.push(res.json('data.transaction_id'));
      }
    });

    sleep(0.2); // Small delay between payments
  }

  sleep(1); // Wait for payments to process

  // Check status of all created transactions
  transactionIds.forEach((transactionId, index) => {
    group(`Transaction Status ${index + 1} - ${transactionId}`, () => {
      const res = http.get(`${BASE_URL}/api/qris/status/${transactionId}`, {
        headers: authHeaders(data.token),
      });
      check(res, {
        'status 200': (r) => r.status === 200,
        'has status': (r) => r.json('data.status') !== undefined,
      });
    });

    sleep(0.1); // Small delay between status checks
  });

  sleep(1);
