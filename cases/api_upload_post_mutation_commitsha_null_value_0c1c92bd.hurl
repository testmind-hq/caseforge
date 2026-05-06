# ── POST /api/upload - mutation: commitSha null value ──
# case_id=TC-0c1c92bd
# case_name=POST /api/upload - mutation: commitSha null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "instance",
  "commitSha": null,
  "service": "sufficient",
  "specContent": "ours"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

