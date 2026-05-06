# ── POST /auth/register - [field_boundary] password invalid_below_min ──
# case_id=TC-29d13f96
# case_name=POST /auth/register - [field_boundary] password invalid_below_min
# step_id=step-main
# step_type=test
# technique=field_boundary
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "carmelmaldonado@schwartz.org",
  "password": "aaaaaaa"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

