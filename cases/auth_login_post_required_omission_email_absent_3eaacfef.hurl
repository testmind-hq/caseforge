# ── POST /auth/login - [required_omission] email absent ──
# case_id=TC-3eaacfef
# case_name=POST /auth/login - [required_omission] email absent
# step_id=step-main
# step_type=test
# technique=required_omission
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "password": "abroad"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

