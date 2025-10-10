import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '30s', target: 200 }, // Ramp-up to 200 users over 30 seconds
        { duration: '1m', target: 200 },  // Stay at 200 users for 1 minute
        { duration: '10s', target: 0 },   // Ramp-down to 0 users over 10 seconds
    ],
};

export default function () {
    // 1. Create user counter
    let username = 'user' + Math.floor(Math.random() * 1000000);
    let createRes = http.post('http://locntp-user.zalopay.vn/api/v1/users', JSON.stringify({ username: username }), {
        headers: {
            'Content-Type': 'application/json',
        },
    });
    check(createRes, { 'create user counter status is 201': (r) => r.status === 201 });
    let userId = createRes.json().data.id;

    // 2. Increase user counter
    let incRes = http.put(`http://locntp-user.zalopay.vn/api/v1/users/${userId}/increment`);
    check(incRes, { 'increase user counter status is 200': (r) => r.status === 200 });

    // 3. Get user counter
    let getRes = http.get(`http://locntp-user.zalopay.vn/api/v1/users/${userId}/count`);
    check(getRes, { 'get user counter status is 200': (r) => r.status === 200 });

    // 4. Delete user counter
    let delRes = http.del(`http://locntp-user.zalopay.vn/api/v1/users/${userId}`);
    check(delRes, { 'delete user counter status is 200': (r) => r.status === 200 });

    sleep(1);
}