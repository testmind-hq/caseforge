# ── POST /api/upload - mutation: commitSha integer instead of string ──
# case_id=TC-b1212f34
# case_name=POST /api/upload - mutation: commitSha integer instead of string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "instance",
  "commitSha": 12345,
  "service": "sufficient",
  "specContent": "ours"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

