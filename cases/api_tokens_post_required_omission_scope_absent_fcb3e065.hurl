# ── POST /api/tokens - [required_omission] scope absent ──
# case_id=TC-fcb3e065
# case_name=POST /api/tokens - [required_omission] scope absent
# step_id=step-main
# step_type=test
# technique=required_omission
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Macey Wolfe"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

