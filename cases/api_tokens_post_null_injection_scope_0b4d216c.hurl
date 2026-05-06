# ── POST /api/tokens - null injection: scope ──
# case_id=TC-0b4d216c
# case_name=POST /api/tokens - null injection: scope
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Evelyn Coleman",
  "scope": null
}
```

HTTP 422

