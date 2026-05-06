# ── POST /api/upload - mutation: branch empty string ──
# case_id=TC-cac690c1
# case_name=POST /api/upload - mutation: branch empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "",
  "commitSha": "heavily",
  "service": "sufficient",
  "specContent": "ours"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

