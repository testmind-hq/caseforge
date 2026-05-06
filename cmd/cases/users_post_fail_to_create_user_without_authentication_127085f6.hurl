# ── Fail to create user without authentication ──
# case_id=TC-127085f6
# case_name=Fail to create user without authentication
# step_id=step-1
# step_type=test
# technique=ask
# priority=P2

POST {{base_url}}/users
```json
{
  "email": "unauth@example.com",
  "name": "Unauthorized User",
  "password": "SecurePass123!"
}
```

HTTP 401

