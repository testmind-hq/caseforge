# ── POST /api/tokens - wrong content-type (text/plain) ──
# case_id=TC-b0b71990
# case_name=POST /api/tokens - wrong content-type (text/plain)
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: text/plain
```json
{
  "name": "Evelyn Coleman",
  "scope": "write"
}
```

HTTP 415

