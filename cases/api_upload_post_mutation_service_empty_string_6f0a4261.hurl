# ── POST /api/upload - mutation: service empty string ──
# case_id=TC-6f0a4261
# case_name=POST /api/upload - mutation: service empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "instance",
  "commitSha": "heavily",
  "service": "",
  "specContent": "ours"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

