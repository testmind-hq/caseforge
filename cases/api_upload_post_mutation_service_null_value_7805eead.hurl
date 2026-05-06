# ── POST /api/upload - mutation: service null value ──
# case_id=TC-7805eead
# case_name=POST /api/upload - mutation: service null value
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
  "service": null,
  "specContent": "ours"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

