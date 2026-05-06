# ══════════════════════════════════════════════════
# sequence chain: /api/admin/teams/{id}/grants → DELETE /api/admin/users/{id}
# case_id=TC-1e93f696
# case_name=sequence chain: /api/admin/teams/{id}/grants → DELETE /api/admin/users/{id}
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
    "next"
  ],
  "expiresAt": "1953-08-22T03:36:54Z",
  "granteeTeamId": "4ec6231f-137f-4153-97d0-8c43294d0bd2",
  "granteeUserId": "94e4e393-307c-46af-870b-f6f1a737e66b",
  "serviceId": "67af3e57-44c9-4422-ae15-53de1e10b9a7"
}
```

HTTP *

[Captures]
id: jsonpath "$.id"

[Asserts]
status < 300

# ── use via DELETE /api/admin/users/{id} [test] ──
# step_id=step-test
# step_type=test
# title=use via DELETE /api/admin/users/{id}
# depends_on=step-setup

DELETE {{base_url}}/api/admin/users/{{id}}

HTTP *

[Asserts]
status < 300

