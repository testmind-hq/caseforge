# ── POST /api/upload - missing required field "specContent" ──
# case_id=TC-fccdadb2
# case_name=POST /api/upload - missing required field "specContent"
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
  "service": "ears"
}
```

HTTP 422

