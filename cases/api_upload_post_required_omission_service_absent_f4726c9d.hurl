# ── POST /api/upload - [required_omission] service absent ──
# case_id=TC-f4726c9d
# case_name=POST /api/upload - [required_omission] service absent
# step_id=step-main
# step_type=test
# technique=required_omission
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "whenever",
  "commitSha": "himself",
  "specContent": "did"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

