import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export let options = {
  stages: [
    { duration: '30s', target: 500 },
    { duration: '1m', target: 1000 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(99)<=50'],
    http_req_failed: ['rate<0.0001'],
  },
};

function getToken() {
  let username = `user_${randomString(8)}`;
  let payload = JSON.stringify({ username: username, password: 'password123' });

  let res = http.post('http://localhost:8080/api/auth', payload, {
    headers: { 'Content-Type': 'application/json' },
  });

  check(res, { 'Auth response is 200': (r) => r.status === 200 });

  let authData = JSON.parse(res.body);
  return { token: authData.token, username: username };
}

export default function () {
  let { token, username } = getToken();
  let headers = { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' };

  let buyRes = http.get('http://localhost:8080/api/buy/t-shirt', { headers });
  check(buyRes, { 'Buy response is 200': (r) => r.status === 200 });

  let infoRes = http.get('http://localhost:8080/api/info', { headers });
  check(infoRes, { 'Info response is 200': (r) => r.status === 200 });

  // let receiver = `user_${randomString(8)}`;
  // http.post('http://localhost:8080/api/auth', JSON.stringify({ username: receiver, password: 'password123' }), { headers });

  // let sendCoinRes = http.post(
  //   'http://localhost:8080/api/sendCoin',
  //   JSON.stringify({ toUser: receiver, amount: 1 }),
  //   { headers }
  // );

  // check(sendCoinRes, { 'SendCoin response is 200': (r) => r.status === 200 });

  sleep(1);
}
