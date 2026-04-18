import http from 'k6/http';
import { check, group, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 100 },
    { duration: '120s', target: 100 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';
const JWT_TOKEN = __ENV.JWT_TOKEN || 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoidXNlcl8xMjMiLCJleHAiOjE3MzUzNjAwMDB9.example'; // Mock JWT
const SAMPLE_QRIS = '00020101021126690021ID.CO.BANKMANDIRI.WWW01189360000801299399930211712993999340303UKE51440014ID.CO.QRIS.WWW0215ID10232756067300303UKE5204274153033605802ID5912M%20Ivan%20Store6015Jakarta%20Timur63045F26';

export default function () {
  let transactionIds = []; // Store transaction IDs from payments

  group('QRIS Inquiry - Normal Network', () => {
    const res = http.get(`${BASE_URL}/api/qris/inquiry/${SAMPLE_QRIS}`, {
      headers: {
        'Authorization': `Bearer ${JWT_TOKEN}`,
        'Content-Type': 'application/json',
      },
    });
    check(res, {
      'status 200': (r) => r.status === 200,
      'has metadata': (r) => r.json('metadata') !== undefined,
      'response time < 500ms': (r) => r.timings.duration < 500,
    });
  });

  sleep(0.5);

  // Create 5 payments and store their transaction IDs
  for (let i = 0; i < 5; i++) {
    group(`QRIS Payment ${i + 1} - Normal Network`, () => {
      const body = JSON.stringify({
        inquiry_id: 'inq_' + Math.random().toString(36).substr(2, 9),
        amount: Math.floor(Math.random() * 100000) + 1000,
        payment_method: 'balance',
        pincode: '123456',
      });
      const res = http.post(`${BASE_URL}/api/qris/payment`, body, {
        headers: {
          'Authorization': `Bearer ${JWT_TOKEN}`,
          'Content-Type': 'application/json',
        },
      });
      check(res, {
        'status 202': (r) => r.status === 202,
        'has transaction_id': (r) => r.json('data.transaction_id') !== undefined,
        'response time < 1000ms': (r) => r.timings.duration < 1000,
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
    group(`Transaction Status ${index + 1} - Normal Network - ${transactionId}`, () => {
      const res = http.get(`${BASE_URL}/api/qris/status/${transactionId}`, {
        headers: {
          'Authorization': `Bearer ${JWT_TOKEN}`,
          'Content-Type': 'application/json',
        },
      });
      check(res, {
        'status 200': (r) => r.status === 200,
        'has status': (r) => r.json('data.status') !== undefined,
        'response time < 500ms': (r) => r.timings.duration < 500,
      });
    });

    sleep(0.1); // Small delay between status checks
  });

  sleep(1);