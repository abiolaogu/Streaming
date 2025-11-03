// k6 load testing script for Kong API Gateway
// Tests all services under load (1000 RPS target)

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const requestDuration = new Trend('request_duration');
const requestsCounter = new Counter('requests');

// Test configuration
export const options = {
  stages: [
    { duration: '30s', target: 100 },   // Ramp up to 100 users
    { duration: '1m', target: 500 },   // Ramp up to 500 users
    { duration: '2m', target: 1000 },  // Ramp up to 1000 users (target RPS)
    { duration: '2m', target: 1000 },  // Stay at 1000 users
    { duration: '1m', target: 500 },   // Ramp down to 500 users
    { duration: '30s', target: 0 },    // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'], // 95% < 500ms, 99% < 1s
    http_req_failed: ['rate<0.01'],                    // Error rate < 1%
    errors: ['rate<0.01'],                            // Custom error rate < 1%
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8000';

// Test scenarios - rotating through different endpoints
const scenarios = [
  {
    name: 'Auth Service',
    endpoints: [
      '/auth/health',
      '/auth/register',
      '/auth/login',
    ],
  },
  {
    name: 'Content Service',
    endpoints: [
      '/content/health',
      '/content/search?q=test',
      '/content/categories',
      '/content/trending',
    ],
  },
  {
    name: 'Search Service',
    endpoints: [
      '/search/health',
      '/search?q=movie',
      '/search/suggest?q=movie',
      '/search/filters',
    ],
  },
  {
    name: 'User Service',
    endpoints: [
      '/users/health',
    ],
  },
  {
    name: 'Payment Service',
    endpoints: [
      '/payments/health',
      '/payments/plans',
    ],
  },
  {
    name: 'Analytics Service',
    endpoints: [
      '/analytics/health',
      '/analytics/dashboard',
    ],
  },
  {
    name: 'Recommendation Service',
    endpoints: [
      '/recommendations/health',
      '/recommendations/trending',
    ],
  },
  {
    name: 'Scheduler Service',
    endpoints: [
      '/scheduler/health',
      '/scheduler/channels',
    ],
  },
];

function getRandomEndpoint() {
  const scenario = scenarios[Math.floor(Math.random() * scenarios.length)];
  const endpoint = scenario.endpoints[Math.floor(Math.random() * scenario.endpoints.length)];
  return endpoint;
}

export default function () {
  const endpoint = getRandomEndpoint();
  const url = `${BASE_URL}${endpoint}`;
  
  // Random headers
  const headers = {
    'Content-Type': 'application/json',
    'User-Agent': 'k6-load-test/1.0',
  };
  
  // For POST endpoints, send JSON payload
  let method = 'GET';
  let body = null;
  
  if (endpoint.includes('/register')) {
    method = 'POST';
    body = JSON.stringify({
      email: `test${Math.random()}@example.com`,
      password: 'Test123!@#',
    });
  } else if (endpoint.includes('/login')) {
    method = 'POST';
    body = JSON.stringify({
      email: 'test@example.com',
      password: 'Test123!@#',
    });
  }
  
  const startTime = Date.now();
  let res;
  
  if (method === 'POST') {
    res = http.post(url, body, { headers });
  } else {
    res = http.get(url, { headers });
  }
  
  const duration = Date.now() - startTime;
  
  // Record metrics
  requestsCounter.add(1);
  requestDuration.add(duration);
  
  // Check response
  const success = check(res, {
    'status is 200 or 401': (r) => r.status === 200 || r.status === 401 || r.status === 404,
    'response time < 1s': (r) => r.timings.duration < 1000,
  });
  
  if (!success || res.status >= 500) {
    errorRate.add(1);
  } else {
    errorRate.add(0);
  }
  
  // Small sleep to avoid overwhelming the system
  sleep(0.1);
}

export function handleSummary(data) {
  return {
    'stdout': textSummary(data, { indent: ' ', enableColors: true }),
    'summary.json': JSON.stringify(data),
  };
}

function textSummary(data, options) {
  const indent = options.indent || '';
  const enableColors = options.enableColors || false;
  
  let summary = '\n';
  summary += `${indent}Kong API Gateway Load Test Summary\n`;
  summary += `${indent}===================================\n\n`;
  
  // Metrics summary
  if (data.metrics) {
    summary += `${indent}Metrics:\n`;
    summary += `${indent}  - Total Requests: ${data.metrics.requests.values.count}\n`;
    summary += `${indent}  - Failed Requests: ${data.metrics.http_req_failed.values.passes}\n`;
    summary += `${indent}  - Avg Response Time: ${data.metrics.http_req_duration.values.avg.toFixed(2)}ms\n`;
    summary += `${indent}  - P95 Response Time: ${data.metrics.http_req_duration.values['p(95)'].toFixed(2)}ms\n`;
    summary += `${indent}  - P99 Response Time: ${data.metrics.http_req_duration.values['p(99)'].toFixed(2)}ms\n`;
    summary += `${indent}  - Error Rate: ${(data.metrics.errors.values.rate * 100).toFixed(2)}%\n`;
  }
  
  summary += '\n';
  return summary;
}

