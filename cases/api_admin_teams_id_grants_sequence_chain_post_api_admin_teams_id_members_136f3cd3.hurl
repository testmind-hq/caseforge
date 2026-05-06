# ══════════════════════════════════════════════════
# sequence chain: /api/admin/teams/{id}/grants → POST /api/admin/teams/{id}/members
# case_id=TC-136f3cd3
# case_name=sequence chain: /api/admin/teams/{id}/grants → POST /api/admin/teams/{id}/members
# case_kind=chain
# priority=P1
# ══════════════════════════════════════════════════

# ── create via POST /api/admin/teams/{id}/grants [setup] ──
# step_id=step-setup
# step_type=setup
# title=create via POST /api/admin/teams/{id}/grants

POST {{base_url}}/api/admin/teams/{id}/grants
Content-Type: application/json
```json
{
  "branches": [
    "wow"
  ],
  "expiresAt": "1972-07-06T21:33:45Z",
  "granteeTeamId": "b14431ac-e726-45f0-93de-31b938772976",
  "granteeUserId": "4d5d2551-5245-4b9f-96e5-0b702e93eff2",
  "serviceId": "fa586d52-80ed-493e-8e6d-6047b31e41fa"
}
```

HTTP *

[Captures]
id: jsonpath "$.id"

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
  "role": "member",
  "userId": "1dd37e1e-0598-4a14-9118-1e52865101d3"
}
```

HTTP *

[Asserts]
status < 300

