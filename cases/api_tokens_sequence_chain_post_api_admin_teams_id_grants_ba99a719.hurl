# ══════════════════════════════════════════════════
# sequence chain: /api/tokens → POST /api/admin/teams/{id}/grants
# case_id=TC-ba99a719
# case_name=sequence chain: /api/tokens → POST /api/admin/teams/{id}/grants
# case_kind=chain
# priority=P1
# ══════════════════════════════════════════════════

# ── create via POST /api/tokens [setup] ─────
# step_id=step-setup
# step_type=setup
# title=create via POST /api/tokens

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Aric Carpenter",
  "scope": "write"
}
```

HTTP *

[Captures]
id: jsonpath "$.id"

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
    "consequence"
  ],
  "expiresAt": "1923-07-31T23:48:34Z",
  "granteeTeamId": "951d9915-63f4-46d3-b5d5-8b170b457b9e",
  "granteeUserId": "bbc3acfe-6b9e-4c9c-bf24-b4d09f78276d",
  "serviceId": "47af9d4e-ddf7-4f73-8a33-2c60da4c1f72"
}
```

HTTP *

[Asserts]
status < 300

