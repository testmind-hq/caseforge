# ══════════════════════════════════════════════════
# sequence chain: /auth/login → POST /api/admin/teams/{id}/grants
# case_id=TC-ba58927e
# case_name=sequence chain: /auth/login → POST /api/admin/teams/{id}/grants
# case_kind=chain
# priority=P1
# ══════════════════════════════════════════════════

# ── create via POST /auth/login [setup] ─────
# step_id=step-setup
# step_type=setup
# title=create via POST /auth/login

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "chaimbird@peters.info",
  "password": "have"
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
    "anybody"
  ],
  "expiresAt": "1900-01-23T02:22:54Z",
  "granteeTeamId": "2c916244-ec7b-46c4-8a46-75d8003b66f2",
  "granteeUserId": "c582e301-b02e-418f-9960-f865b66da97f",
  "serviceId": "eaa19ebb-002b-497c-a98a-0293aa5606ad"
}
```

HTTP *

[Asserts]
status < 300

