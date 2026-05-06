# ── POST /auth/login - valid request with all required fields ──
# case_id=TC-486e8c2a
# case_name=POST /auth/login - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "ezrahowell@franklin.biz",
  "password": "work"
}
```

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.role" exists
jsonpath "$.token" exists
jsonpath "$.userId" exists

