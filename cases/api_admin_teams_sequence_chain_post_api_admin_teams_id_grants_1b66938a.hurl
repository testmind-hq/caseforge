# ══════════════════════════════════════════════════
# sequence chain: /api/admin/teams → POST /api/admin/teams/{id}/grants
# case_id=TC-1b66938a
# case_name=sequence chain: /api/admin/teams → POST /api/admin/teams/{id}/grants
# case_kind=chain
# priority=P1
# ══════════════════════════════════════════════════

# ── create via POST /api/admin/teams [setup] ──
# step_id=step-setup
# step_type=setup
# title=create via POST /api/admin/teams

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Evenings in Oakland invite quieter man.",
  "displayName": "which",
  "name": "Clifton Shields"
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
    "it"
  ],
  "expiresAt": "2001-12-10T08:50:19Z",
  "granteeTeamId": "722fd61c-8b80-44f6-9e81-c9c8550ab73d",
  "granteeUserId": "a1efd1eb-3a36-4f78-85fb-7edd1d4af481",
  "serviceId": "2a7ed0b1-582d-4271-9b40-91828aded5f0"
}
```

HTTP *

[Asserts]
status < 300

