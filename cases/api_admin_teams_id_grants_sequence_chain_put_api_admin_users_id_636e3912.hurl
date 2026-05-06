# ══════════════════════════════════════════════════
# sequence chain: /api/admin/teams/{id}/grants → PUT /api/admin/users/{id}
# case_id=TC-636e3912
# case_name=sequence chain: /api/admin/teams/{id}/grants → PUT /api/admin/users/{id}
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
    "half"
  ],
  "expiresAt": "1911-12-23T17:30:07Z",
  "granteeTeamId": "e275d7a1-f1f0-449b-9962-e43b92698249",
  "granteeUserId": "5a22025f-d28e-4434-9b1d-93bf353fbdb9",
  "serviceId": "71bbc723-acdf-4be2-b56f-e471f9077cc5"
}
```

HTTP *

[Captures]
id: jsonpath "$.id"

[Asserts]
status < 300

# ── use via PUT /api/admin/users/{id} [test] ──
# step_id=step-test
# step_type=test
# title=use via PUT /api/admin/users/{id}
# depends_on=step-setup

PUT {{base_url}}/api/admin/users/{{id}}
Content-Type: application/json
```json
{
  "isActive": true,
  "role": "team_member"
}
```

HTTP *

[Asserts]
status < 300

