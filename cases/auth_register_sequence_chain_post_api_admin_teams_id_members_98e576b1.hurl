# ══════════════════════════════════════════════════
# sequence chain: /auth/register → POST /api/admin/teams/{id}/members
# case_id=TC-98e576b1
# case_name=sequence chain: /auth/register → POST /api/admin/teams/{id}/members
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
  "email": "lukasvalencia@cummings.name",
  "password": "couple"
}
```

HTTP *

[Captures]
id: jsonpath "$.userId"

[Asserts]
status < 300

# ── use via POST /api/admin/teams/{id}/members [test] ──
# step_id=step-test
# step_type=test
# title=use via POST /api/admin/teams/{id}/members
# depends_on=step-setup

POST {{base_url}}/api/admin/teams/{{id}}/members
Content-Type: application/json
```json
{
  "role": "owner",
  "userId": "204452b4-832e-4601-a227-8ecf3cc125ec"
}
```

HTTP *

[Asserts]
status < 300

