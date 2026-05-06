# ── POST /api/upload - invalid branch: empty string violates minLength 1 ──
# case_id=TC-5eb7446c
# case_name=POST /api/upload - invalid branch: empty string violates minLength 1
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "",
  "commitSha": "pack",
  "service": "ears",
  "specContent": "now"
}
```

HTTP 422

