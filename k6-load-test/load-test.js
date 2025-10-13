import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '5m', target: 100 }, // Ramp-up to 100 users over 5 minutes
        { duration: '5m', target: 300 }, // Ramp-up to 300 users over 5 minutes
        { duration: '5m', target: 600 }, // Stay at 600 users for 5 minutes
        { duration: '5m', target: 800 },  // Stay at 800 users for 5 minutes
        { duration: '3m', target: 200 }, // Ramp-down to 200 users over 3 minutes
        { duration: '2m', target: 0 },   // Ramp-down to 0 users over 2 minutes
    ],
};

// random generate username function
function generateUsername() {
    // Combine timestamp + random 3-digit number + worker ID
    // This gives ~10^12 unique possibilities
    const timestamp = Date.now().toString(36); // time component (base36 compress)
    const random = Math.floor(Math.random() * 1e9).toString(36); // random part
    return `user_${timestamp}_${random}`;
}


export default function () {
    // 1. Create user counter
    let username = generateUsername();
    let createRes = http.post('http://locntp-user.zalopay.vn/api/v1/users', JSON.stringify({ username: username }), {
        headers: {
            'Content-Type': 'application/json',
        },
    });
    check(createRes, { 'create user counter status is 201': (r) => r.status === 201 });
    if (createRes.status !== 201) {
        // print error and return if user creation fails
        console.error('User creation failed:', createRes.body);
        return;
    }

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