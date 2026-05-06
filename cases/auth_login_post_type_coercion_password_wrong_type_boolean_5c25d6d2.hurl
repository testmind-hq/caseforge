# ── POST /auth/login - [type_coercion] password wrong_type_boolean ──
# case_id=TC-5c25d6d2
# case_name=POST /auth/login - [type_coercion] password wrong_type_boolean
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "lunasaunders@greene.net",
  "password": true
}
```

HTTP 422

