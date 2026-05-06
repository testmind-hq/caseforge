# ── POST /api/upload - [required_omission] specContent absent ──
# case_id=TC-196e600f
# case_name=POST /api/upload - [required_omission] specContent absent
# step_id=step-main
# step_type=test
# technique=required_omission
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "now",
  "commitSha": "occasionally",
  "service": "might"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

