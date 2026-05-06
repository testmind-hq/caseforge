# ── POST /api/upload - null injection: specContent ──
# case_id=TC-fef2ed50
# case_name=POST /api/upload - null injection: specContent
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "rather",
  "commitSha": "troop",
  "service": "we",
  "specContent": null
}
```

HTTP 422

