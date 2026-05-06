# ── POST /api/upload - missing required field "branch" ──
# case_id=TC-33947120
# case_name=POST /api/upload - missing required field "branch"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "commitSha": "pack",
  "service": "ears",
  "specContent": "now"
}
```

HTTP 422

