# ── POST /api/tokens - null injection: name ──
# case_id=TC-97bd0c77
# case_name=POST /api/tokens - null injection: name
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": null,
  "scope": "write"
}
```

HTTP 422

