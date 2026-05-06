# ══════════════════════════════════════════════════
# sequence chain: /api/admin/teams/{id}/grants → PUT /api/admin/services/{serviceId}/team
# case_id=TC-cafaccf6
# case_name=sequence chain: /api/admin/teams/{id}/grants → PUT /api/admin/services/{serviceId}/team
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
    "am"
  ],
  "expiresAt": "1930-06-02T07:33:10Z",
  "granteeTeamId": "6eb082a3-7a81-4673-b080-6f876150d238",
  "granteeUserId": "9c8b45fd-f191-4a4d-80fd-b8dad10d176a",
  "serviceId": "d078acf6-4a9a-463a-9632-1d93b5a7ecfa"
}
```

HTTP *

[Captures]
serviceId: jsonpath "$.id"

[Asserts]
status < 300

# ── use via PUT /api/admin/services/{serviceId}/team [test] ──
# step_id=step-test
# step_type=test
# title=use via PUT /api/admin/services/{serviceId}/team
# depends_on=step-setup

PUT {{base_url}}/api/admin/services/{{serviceId}}/team
Content-Type: application/json
```json
{
  "teamId": "ef302aa8-fd8d-4fd6-9798-6d57d88f7ac6"
}
```

HTTP *

[Asserts]
status < 300

