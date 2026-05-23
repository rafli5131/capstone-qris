import http from 'k6/http';

const TEST_PASSWORD = 'password123';

function uniqueUsername(prefix) {
  const random = Math.random().toString(36).slice(2, 10);
  return `${prefix}${Date.now()}${random}`;
}

export function registerUser(baseUrl, prefix) {
  const body = JSON.stringify({
    username: uniqueUsername(prefix),
    password: TEST_PASSWORD,
    initial_balance: 1000000,
  });

  const res = http.post(`${baseUrl}/api/auth/register`, body, {
    tags: { endpoint: 'auth_register' },
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
  });

  if (res.status !== 201) {
    throw new Error(`${prefix} register failed: ${res.status} ${res.body}`);
  }

  const token = res.json('data.token');
  if (!token) {
    throw new Error(`${prefix} register response did not include token`);
  }

  return token;
}

export function authHeaders(token) {
  return {
    Accept: 'application/json',
    Authorization: `Bearer ${token}`,
  };
}
