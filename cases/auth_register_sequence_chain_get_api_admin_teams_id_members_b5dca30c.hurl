# ══════════════════════════════════════════════════
# sequence chain: /auth/register → GET /api/admin/teams/{id}/members
# case_id=TC-b5dca30c
# case_name=sequence chain: /auth/register → GET /api/admin/teams/{id}/members
# case_kind=chain
# priority=P1
# ══════════════════════════════════════════════════

# ── create via POST /auth/register [setup] ──
# step_id=step-setup
# step_type=setup
# title=create via POST /auth/register

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "greggburns@spencer.info",
  "password": "motivation"
}
```

HTTP *

[Captures]
id: jsonpath "$.userId"

[Asserts]
status < 300

# ── use via GET /api/admin/teams/{id}/members [test] ──
# step_id=step-test
# step_type=test
# title=use via GET /api/admin/teams/{id}/members
# depends_on=step-setup

GET {{base_url}}/api/admin/teams/{{id}}/members

HTTP *

[Asserts]
status < 300

