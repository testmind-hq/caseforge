# ── POST /auth/register - valid request with all required fields ──
# case_id=TC-787a33be
# case_name=POST /auth/register - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "alessandravaldez@daniel.net",
  "password": "who"
}
```

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.userId" exists

