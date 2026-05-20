import http from 'k6/http';
import { check, group, sleep } from 'k6';
import crypto from 'k6/crypto';

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
const CLIENT_KEY = 'MK-9921-X';
const CLIENT_SECRET = 'super-secret-key-change-in-production';
const SAMPLE_QRIS = '00020101021126690021ID.CO.BANKMANDIRI.WWW01189360000801299399930211712993999340303UKE51440014ID.CO.QRIS.WWW0215ID10232756067300303UKE5204274153033605802ID5912M%20Ivan%20Store6015Jakarta%20Timur63045F26';

function generateSignature(payload) {
  const encoder = new TextEncoder();
  const data = encoder.encode(payload);
  const secretBytes = encoder.encode(CLIENT_SECRET);
  const hmac = crypto.hmac('sha256', secretBytes, data);
  return crypto.hex(hmac);
}

function getTimestamp() {
  return new Date().toISOString();
}

export default function () {
  const timestamp = getTimestamp();

  group('QRIS Inquiry', () => {
    const sig = generateSignature('');
    const res = http.get(`${BASE_URL}/api/qris/inquiry/${SAMPLE_QRIS}`, {
      headers: {
        'X-Timestamp': timestamp,
        'X-Client-Key': CLIENT_KEY,
        'X-Signature': sig,
        'Content-Type': 'application/json',
      },
    });
    check(res, {
      'status 200': (r) => r.status === 200,
      'has metadata': (r) => r.json('metadata') !== undefined,
    });
  });

  sleep(0.5);

  group('QRIS Payment', () => {
    const body = JSON.stringify({
      inquiry_id: 'inq_12345678',
      user_id: 'user_123',
      amount: 50000,
      payment_method: 'balance',
      pincode: '123456',
    });
    const sig = generateSignature(body);
    const res = http.post(`${BASE_URL}/api/qris/payment`, body, {
      headers: {
        'X-Timestamp': timestamp,
        'X-Client-Key': CLIENT_KEY,
        'X-Signature': sig,
        'Content-Type': 'application/json',
      },
    });
    check(res, {
      'status 202': (r) => r.status === 202,
      'has transaction_id': (r) => r.json('data.transaction_id') !== undefined,
    });
  });

  sleep(1);
}
