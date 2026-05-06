# ── POST /auth/register - [required_omission] email absent ──
# case_id=TC-b724df31
# case_name=POST /auth/register - [required_omission] email absent
# step_id=step-main
# step_type=test
# technique=required_omission
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "password": "themselves"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

