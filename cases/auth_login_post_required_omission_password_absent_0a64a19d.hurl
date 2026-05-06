# ── POST /auth/login - [required_omission] password absent ──
# case_id=TC-0a64a19d
# case_name=POST /auth/login - [required_omission] password absent
# step_id=step-main
# step_type=test
# technique=required_omission
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "darylfarrell@santiago.org"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

