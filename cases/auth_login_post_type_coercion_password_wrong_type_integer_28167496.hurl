# ── POST /auth/login - [type_coercion] password wrong_type_integer ──
# case_id=TC-28167496
# case_name=POST /auth/login - [type_coercion] password wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "lunasaunders@greene.net",
  "password": 123
}
```

HTTP 422

