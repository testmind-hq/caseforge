# ── POST /api/upload - mutation: branch integer instead of string ──
# case_id=TC-416a96c1
# case_name=POST /api/upload - mutation: branch integer instead of string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": 12345,
  "commitSha": "heavily",
  "service": "sufficient",
  "specContent": "ours"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

