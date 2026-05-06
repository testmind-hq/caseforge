# ── POST /api/upload - invalid service: empty string violates minLength 1 ──
# case_id=TC-8389dd21
# case_name=POST /api/upload - invalid service: empty string violates minLength 1
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "recline",
  "commitSha": "pack",
  "service": "",
  "specContent": "now"
}
```

HTTP 422

