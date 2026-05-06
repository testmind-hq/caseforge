# ── POST /api/upload - mutation: branch null value ──
# case_id=TC-9f510ed7
# case_name=POST /api/upload - mutation: branch null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": null,
  "commitSha": "heavily",
  "service": "sufficient",
  "specContent": "ours"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

