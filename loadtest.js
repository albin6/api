import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '30s', target: 50 },
        { duration: '1m', target: 50 },
        { duration: '30s', target: 0 },
    ],
    thresholds: {
        http_req_duration: ['p(95)<200'],
    },
};

export default function () {
    let res = http.get('http://localhost:8080/students?search=albin&page=1&limit=10');
    check(res, { 'status was 200': (r) => r.status == 200 });
    sleep(1);
}
