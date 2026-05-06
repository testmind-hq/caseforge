# ══════════════════════════════════════════════════
# sequence chain: /auth/login → GET /api/admin/teams/{id}/grants
# case_id=TC-4f853ed4
# case_name=sequence chain: /auth/login → GET /api/admin/teams/{id}/grants
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
  "email": "frankiewebb@davies.org",
  "password": "for"
}
```

HTTP *

[Captures]
id: jsonpath "$.userId"

[Asserts]
status < 300

# ── use via GET /api/admin/teams/{id}/grants [test] ──
# step_id=step-test
# step_type=test
# title=use via GET /api/admin/teams/{id}/grants
# depends_on=step-setup

GET {{base_url}}/api/admin/teams/{{id}}/grants

HTTP *

[Asserts]
status < 300

