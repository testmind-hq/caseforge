# ── POST /auth/register - [required_omission] password absent ──
# case_id=TC-3d6d9a7d
# case_name=POST /auth/register - [required_omission] password absent
# step_id=step-main
# step_type=test
# technique=required_omission
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "artperkins@smith.net"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

