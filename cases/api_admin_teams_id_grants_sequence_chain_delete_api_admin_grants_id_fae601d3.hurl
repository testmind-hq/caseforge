# ══════════════════════════════════════════════════
# sequence chain: /api/admin/teams/{id}/grants → DELETE /api/admin/grants/{id}
# case_id=TC-fae601d3
# case_name=sequence chain: /api/admin/teams/{id}/grants → DELETE /api/admin/grants/{id}
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
    "in"
  ],
  "expiresAt": "1934-04-27T17:54:54Z",
  "granteeTeamId": "ef7ba0e3-e654-4cbe-a8db-7d80ae34554a",
  "granteeUserId": "6b8cf351-2a07-4e9b-af8d-93adadf31af4",
  "serviceId": "4af3c971-e3ff-4038-8eec-7562f600ef7e"
}
```

HTTP *

[Captures]
id: jsonpath "$.id"

[Asserts]
status < 300

# ── use via DELETE /api/admin/grants/{id} [test] ──
# step_id=step-test
# step_type=test
# title=use via DELETE /api/admin/grants/{id}
# depends_on=step-setup

DELETE {{base_url}}/api/admin/grants/{{id}}

HTTP *

[Asserts]
status < 300

