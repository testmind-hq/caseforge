# ── POST /api/tokens - missing required field "scope" ──
# case_id=TC-aa18d499
# case_name=POST /api/tokens - missing required field "scope"
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P1

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Lawrence Braun"
}
```

HTTP 422

