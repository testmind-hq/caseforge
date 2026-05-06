# ── POST /api/admin/teams/{id}/members - mutation: role empty string ──
# case_id=TC-0cb69d90
# case_name=POST /api/admin/teams/{id}/members - mutation: role empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "",
  "userId": "eb5af601-571e-49ce-a28d-f33fe87bc344"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

