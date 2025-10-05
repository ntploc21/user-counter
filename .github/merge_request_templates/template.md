## Summary:
<!-- Provide a brief overview of what the merge request does. Why are these changes necessary? -->
- What problem is this solving?
- What changes were made in this, MR?

## Implementation Details:
<!-- Provide a more in-depth explanation of how the problem was solved. Include key technical decisions or trade-offs. -->
<!-- - Design choices and reasoning:-->
<!--     - [Example: "Switched to a concurrent model to improve performance by X%."]-->
<!--     - [Mention new/updated architectural patterns if applicable]-->
<!-- - Key files or components impacted:-->
<!--     - `pkg/controller/mycontroller.go`-->
<!--     - `internal/utils/`-->

## **How to Test:**

### 1. API 1

**Example:**
- **Endpoint:** `POST /api/v1/users/login`
- **Request:**
    ```json
    {
        "username": "test_user",
        "password": "password123"
    }
    ```
- **Expected Response (200 OK):**
    ```json
    {
        "status": "success",
        "data": {
            "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
            "user": {
                "id": "12345",
                "username": "test_user"
            }
        }
    }
    ```
- **Steps to Verify:**
    1. Response status code should be `200 OK`.
    2. Ensure a valid `token` is present.
    3. Confirm the `user` object matches expected data.

- **Sample cURL Command:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/users/login \
       -H "Content-Type: application/json" \
       -d '{"username": "test_user", "password": "password123"}'
   ```

## Deployment Notes:
<!-- Mention if special deployment or migration steps are needed, such as running migrations, scaling changes, or feature flag toggling. -->
- [ ] Any database migrations?
- [ ] Requires feature flag?
- [ ] Special configurations? (env variables, secrets, etc.)

## Code Quality and Guidelines Checklist:
- [ ] MR type: `Feature`/`Bug`/`Improvement`
- [ ] Code follows the project's style guidelines
- [ ] Code is self-documented or has appropriate comments
- [ ] New dependencies are well-justified and necessary
- [ ] Error handling is accounted for (no ignored errors)
- [ ] No major performance bottlenecks introduced
- [ ] API documentation updated (if applicable)
- [ ] Does the MR change API endpoints or data formats?

## Reviewers Checklist:
- [ ] Functionality: Does the code do what it is intended to do?
- [ ] Readability: Is the code easy to understand?
- [ ] Tests: Are tests comprehensive and do they cover edge cases?
- [ ] Performance: Is the code efficient and scalable?
- [ ] Security: Are there any potential security risks introduced by this MR?