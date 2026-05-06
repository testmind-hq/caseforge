# ── POST /api/upload - mutation: branch oversized string (300 chars) ──
# case_id=TC-75d60dab
# case_name=POST /api/upload - mutation: branch oversized string (300 chars)
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "commitSha": "heavily",
  "service": "sufficient",
  "specContent": "ours"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

