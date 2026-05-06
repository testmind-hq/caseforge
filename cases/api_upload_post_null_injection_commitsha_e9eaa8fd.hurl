# ── POST /api/upload - null injection: commitSha ──
# case_id=TC-e9eaa8fd
# case_name=POST /api/upload - null injection: commitSha
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "rather",
  "commitSha": null,
  "service": "we",
  "specContent": "usually"
}
```

HTTP 422

