# ── POST /auth/login - missing required field "password" ──
# case_id=TC-a6bbbeb7
# case_name=POST /auth/login - missing required field "password"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "ebonysilva@mendez.info"
}
```

HTTP 422

