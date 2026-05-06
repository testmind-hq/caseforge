# ── POST /api/upload - wrong content-type (text/plain) ──
# case_id=TC-863dd501
# case_name=POST /api/upload - wrong content-type (text/plain)
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/api/upload
Content-Type: text/plain
```json
{
  "branch": "rather",
  "commitSha": "troop",
  "service": "we",
  "specContent": "usually"
}
```

HTTP 415

