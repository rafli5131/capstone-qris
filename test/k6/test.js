import http from 'k6/http';
import { check, group, sleep } from 'k6';

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
const CLIENT_ID = 'MK-9921-X';
const CLIENT_SECRET = 'penyakit-capstone-gila';
const SAMPLE_QRIS = '00020101021126690021ID.CO.BANKMANDIRI.WWW01189360000801299399930211712993999340303UKE51440014ID.CO.QRIS.WWW0215ID10232756067300303UKE5204274153033605802ID5912M%20Ivan%20Store6015Jakarta%20Timur63045F26';


export default function () {

  group('QRIS Inquiry', () => {
    const sig = generateSignature('');
    const res = http.get(`${BASE_URL}/api/qris/inquiry/${SAMPLE_QRIS}`, {
      headers: {
        'X-Client-ID': CLIENT_ID,
        'X-Client-Key': CLIENT_SECRET,
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
      amount: 1000,
      payment_method: 'balance',
      pincode: '123456',
    });
    const sig = generateSignature(body);
    const res = http.post(`${BASE_URL}/api/qris/payment`, body, {
      headers: {
        'X-Client-ID': CLIENT_ID,
        'X-Client-Key': CLIENT_SECRET,
        'Content-Type': 'application/json',
      },
    });
    check(res, {
      'status 202': (r) => r.status === 202,
      'has transaction_id': (r) => r.json('data.transaction_id') !== undefined,
    });
  });

  group('Transaction Status', () => {
    const transactionID = 'tx_12345678';
    const res = http.get(`${BASE_URL}/api/transaction/status/${transactionID}`, {
      headers: {
        'X-Client-ID': CLIENT_ID,
        'X-Client-Key': CLIENT_SECRET,
        'Content-Type': 'application/json',
      },
    });
    check(res, {
      'status 200': (r) => r.status === 200,
      'has status': (r) => r.json('data.status') !== undefined,
    });
  });

  sleep(1);
}
