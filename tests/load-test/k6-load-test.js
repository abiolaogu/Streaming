import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const requestDuration = new Trend('request_duration');

export const options = {
  stages: [
    { duration: '2m', target: 100 },  // Ramp up to 100 users
    { duration: '5m', target: 100 },    // Stay at 100 users
    { duration: '2m', target: 200 },    // Ramp up to 200 users
    { duration: '5m', target: 200 },    // Stay at 200 users
    { duration: '2m', target: 300 },    // Ramp up to 300 users
    { duration: '5m', target: 300 },    // Stay at 300 users
    { duration: '2m', target: 0 },      // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'],
    http_req_failed: ['rate<0.01'],
    errors: ['rate<0.01'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://api.streamverse.io';
const AUTH_TOKEN = __ENV.AUTH_TOKEN || '';

export default function () {
  // Auth Service
  const loginRes = http.post(`${BASE_URL}/auth/login`, JSON.stringify({
    email: `user${Math.floor(Math.random() * 1000)}@test.com`,
    password: 'test123',
  }), {
    headers: { 'Content-Type': 'application/json' },
  });
  
  const loginOk = check(loginRes, {
    'login status is 200 or 401': (r) => r.status === 200 || r.status === 401,
  });
  errorRate.add(!loginOk);
  
  let token = AUTH_TOKEN;
  if (loginRes.status === 200) {
    token = JSON.parse(loginRes.body).token;
  }
  
  // Content Service
  const contentRes = http.get(`${BASE_URL}/content/trending`, {
    headers: { 'Authorization': `Bearer ${token}` },
  });
  
  const contentOk = check(contentRes, {
    'content status is 200': (r) => r.status === 200,
  });
  errorRate.add(!contentOk);
  
  // Search Service
  const searchRes = http.get(`${BASE_URL}/search?q=action`, {
    headers: { 'Authorization': `Bearer ${token}` },
  });
  
  const searchOk = check(searchRes, {
    'search status is 200': (r) => r.status === 200,
  });
  errorRate.add(!searchOk);
  
  // Streaming Service
  if (contentRes.status === 200) {
    const content = JSON.parse(contentRes.body);
    if (content.items && content.items.length > 0) {
      const contentId = content.items[0].id;
      const tokenRes = http.post(`${BASE_URL}/streaming/token`, JSON.stringify({
        content_id: contentId,
      }), {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });
      
      check(tokenRes, {
        'token status is 200': (r) => r.status === 200,
      });
    }
  }
  
  requestDuration.add(contentRes.timings.duration);
  
  sleep(1);
}

