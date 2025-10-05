import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    vus: 10, // number of virtual users
    duration: '30s', // test duration
};

export default function () {
    // 1. Create user counter
    let username = 'user' + Math.floor(Math.random() * 1000000);
    let createRes = http.post('http://localhost:63061/api/v1/users', JSON.stringify({ username: username }), {
        headers: {
            'Content-Type': 'application/json',
        },
    });
    check(createRes, { 'create user counter status is 201': (r) => r.status === 201 });
    let userId = createRes.json().data.id;

    // 2. Increase user counter
    let incRes = http.put(`http://localhost:63061/api/v1/users/${userId}/increment`);
    check(incRes, { 'increase user counter status is 200': (r) => r.status === 200 });

    // 3. Get user counter
    let getRes = http.get(`http://localhost:63061/api/v1/users/${userId}/count`);
    check(getRes, { 'get user counter status is 200': (r) => r.status === 200 });

    sleep(1);
}