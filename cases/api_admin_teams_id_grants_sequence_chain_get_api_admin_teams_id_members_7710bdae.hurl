# ══════════════════════════════════════════════════
# sequence chain: /api/admin/teams/{id}/grants → GET /api/admin/teams/{id}/members
# case_id=TC-7710bdae
# case_name=sequence chain: /api/admin/teams/{id}/grants → GET /api/admin/teams/{id}/members
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
    "place"
  ],
  "expiresAt": "1973-01-05T11:42:04Z",
  "granteeTeamId": "58c7d788-061b-4021-9e8c-01942f155464",
  "granteeUserId": "1b70dc76-c2d3-4e62-9f5d-22c8319dc0a2",
  "serviceId": "a31b4938-a01f-4bc1-80fe-f165a18d784e"
}
```

HTTP *

[Captures]
id: jsonpath "$.id"

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

