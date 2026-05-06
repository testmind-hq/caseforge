# ══════════════════════════════════════════════════
# sequence chain: /auth/register → POST /api/admin/teams/{id}/grants
# case_id=TC-10533daf
# case_name=sequence chain: /auth/register → POST /api/admin/teams/{id}/grants
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
  "email": "chaunceyjacobi@white.com",
  "password": "that"
}
```

HTTP *

[Captures]
id: jsonpath "$.userId"

[Asserts]
status < 300

# ── use via POST /api/admin/teams/{id}/grants [test] ──
# step_id=step-test
# step_type=test
# title=use via POST /api/admin/teams/{id}/grants
# depends_on=step-setup

POST {{base_url}}/api/admin/teams/{{id}}/grants
Content-Type: application/json
```json
{
  "branches": [
    "disregard"
  ],
  "expiresAt": "2003-09-24T09:23:31Z",
  "granteeTeamId": "c727d010-3eb5-469f-93d2-a46ab145fcf5",
  "granteeUserId": "9f6fa71f-b14f-4fe8-bd62-fe79743d34db",
  "serviceId": "1f968d6d-ab6e-4d94-b8de-a0df2b4a5209"
}
```

HTTP *

[Asserts]
status < 300

