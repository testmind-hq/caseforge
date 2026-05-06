# ── POST /api/upload - null injection: service ──
# case_id=TC-b8cf0920
# case_name=POST /api/upload - null injection: service
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
  "service": null,
  "specContent": "usually"
}
```

HTTP 422

