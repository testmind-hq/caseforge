# ── POST /api/upload - mutation: commitSha empty string ──
# case_id=TC-f30e852c
# case_name=POST /api/upload - mutation: commitSha empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "instance",
  "commitSha": "",
  "service": "sufficient",
  "specContent": "ours"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

