# ── POST /api/tokens - name at min_minus_one_invalid boundary ──
# case_id=TC-d08f5a90
# case_name=POST /api/tokens - name at min_minus_one_invalid boundary
# step_id=step-main
# step_type=test
# technique=boundary_value
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "e",
  "scope": "read"
}
```

HTTP 422

