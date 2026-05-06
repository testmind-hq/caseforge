# ══════════════════════════════════════════════════
# sequence chain: /auth/login → POST /api/admin/teams/{id}/members
# case_id=TC-b9578186
# case_name=sequence chain: /auth/login → POST /api/admin/teams/{id}/members
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
  "email": "dwightsummers@schuster.org",
  "password": "model"
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
  "userId": "5f656700-5067-4ad1-8384-1fb850bc7bf2"
}
```

HTTP *

[Asserts]
status < 300

