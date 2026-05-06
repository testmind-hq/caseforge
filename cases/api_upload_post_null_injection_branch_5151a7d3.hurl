# ── POST /api/upload - null injection: branch ──
# case_id=TC-5151a7d3
# case_name=POST /api/upload - null injection: branch
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": null,
  "commitSha": "troop",
  "service": "we",
  "specContent": "usually"
}
```

HTTP 422

