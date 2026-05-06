# ── POST /auth/register - password at min_valid boundary ──
# case_id=TC-0f0b429e
# case_name=POST /auth/register - password at min_valid boundary
# step_id=step-main
# step_type=test
# technique=boundary_value
# priority=P1

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "kasandravelazquez@willis.org",
  "password": "htnnilAG"
}
```

HTTP 200

[Asserts]
duration < 2000

