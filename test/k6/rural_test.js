import http from 'k6/http';
import { check, group, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 10 },   // Rural: low concurrent users
    { duration: '120s', target: 10 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<3000'], // Rural: higher latency tolerance
    http_req_failed: ['rate<0.05'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';
const JWT_TOKEN = __ENV.JWT_TOKEN || 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoidXNlcl8xMjMiLCJleHAiOjE3MzUzNjAwMDB9.example'; // Mock JWT
const SAMPLE_QRIS = '00020101021126690021ID.CO.BANKMANDIRI.WWW01189360000801299399930211712993999340303UKE51440014ID.CO.QRIS.WWW0215ID10232756067300303UKE5204274153033605802ID5912M%20Ivan%20Store6015Jakarta%20Timur63045F26';

export default function () {
  let transactionIds = []; // Store transaction IDs from payments

  // Simulate rural network latency
  sleep(Math.random() * 1.5 + 0.5); // 0.5-2s delay

  group('QRIS Inquiry - Rural Network', () => {
    const res = http.get(`${BASE_URL}/api/qris/inquiry/${SAMPLE_QRIS}`, {
      headers: {
        'Authorization': `Bearer ${JWT_TOKEN}`,
        'Content-Type': 'application/json',
      },
    });
    check(res, {
      'status 200': (r) => r.status === 200,
      'has metadata': (r) => r.json('metadata') !== undefined,
      'response time < 3s': (r) => r.timings.duration < 3000,
    });
  });

  sleep(Math.random() * 2 + 1); // Additional rural delay

  // Create 5 payments and store their transaction IDs
  for (let i = 0; i < 5; i++) {
    group(`QRIS Payment ${i + 1} - Rural Network`, () => {
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
        'response time < 5s': (r) => r.timings.duration < 5000,
      });

      // Store transaction ID for later status checks
      if (res.status === 202 && res.json('data.transaction_id')) {
        transactionIds.push(res.json('data.transaction_id'));
      }
    });

    sleep(Math.random() * 1 + 0.5); // Rural delay between payments
  }

  sleep(Math.random() * 3 + 1); // Rural network delay

  // Check status of all created transactions
  transactionIds.forEach((transactionId, index) => {
    group(`Transaction Status ${index + 1} - Rural Network - ${transactionId}`, () => {
      const res = http.get(`${BASE_URL}/api/qris/status/${transactionId}`, {
        headers: {
          'Authorization': `Bearer ${JWT_TOKEN}`,
          'Content-Type': 'application/json',
        },
      });
      check(res, {
        'status 200': (r) => r.status === 200,
        'has status': (r) => r.json('data.status') !== undefined,
        'response time < 3s': (r) => r.timings.duration < 3000,
      });
    });

    sleep(Math.random() * 0.5 + 0.2); // Small rural delay between status checks
  });

  sleep(Math.random() * 2 + 1); // Rural network delay