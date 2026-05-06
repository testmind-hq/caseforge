# ── POST /auth/login - [type_coercion] email wrong_type_boolean ──
# case_id=TC-91a4d98b
# case_name=POST /auth/login - [type_coercion] email wrong_type_boolean
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": true,
  "password": "whole"
}
```

HTTP 422

