# ── POST /api/upload - missing required field "service" ──
# case_id=TC-8f85caae
# case_name=POST /api/upload - missing required field "service"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "recline",
  "commitSha": "pack",
  "specContent": "now"
}
```

HTTP 422

