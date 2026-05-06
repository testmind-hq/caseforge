# ── POST /api/upload - [required_omission] branch absent ──
# case_id=TC-893f33e4
# case_name=POST /api/upload - [required_omission] branch absent
# step_id=step-main
# step_type=test
# technique=required_omission
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "commitSha": "where",
  "service": "though",
  "specContent": "wisp"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

