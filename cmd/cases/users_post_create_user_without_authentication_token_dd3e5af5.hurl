# ── Create user without authentication token ──
# case_id=TC-dd3e5af5
# case_name=Create user without authentication token
# step_id=step-1
# step_type=test
# technique=ask
# priority=P1

POST {{base_url}}/users
```json
{
  "email": "unauth@example.com",
  "name": "Unauthorized User",
  "password": "SecurePass123!"
}
```

HTTP 401

