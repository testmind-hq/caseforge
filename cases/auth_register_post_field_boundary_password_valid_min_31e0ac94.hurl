# ── POST /auth/register - [field_boundary] password valid_min ──
# case_id=TC-31e0ac94
# case_name=POST /auth/register - [field_boundary] password valid_min
# step_id=step-main
# step_type=test
# technique=field_boundary
# priority=P1

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "cedrickhermann@morales.org",
  "password": "aaaaaaaa"
}
```

HTTP *

[Asserts]
status >= 200
status < 300

