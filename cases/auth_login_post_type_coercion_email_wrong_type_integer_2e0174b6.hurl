# ── POST /auth/login - [type_coercion] email wrong_type_integer ──
# case_id=TC-2e0174b6
# case_name=POST /auth/login - [type_coercion] email wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": 123,
  "password": "whole"
}
```

HTTP 422

