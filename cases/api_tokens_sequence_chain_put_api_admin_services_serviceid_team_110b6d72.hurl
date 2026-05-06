# ══════════════════════════════════════════════════
# sequence chain: /api/tokens → PUT /api/admin/services/{serviceId}/team
# case_id=TC-110b6d72
# case_name=sequence chain: /api/tokens → PUT /api/admin/services/{serviceId}/team
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
  "name": "Diego Herman",
  "scope": "read"
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
  "teamId": "9e4f4d0e-d5d7-447e-830c-1c638616ddbf"
}
```

HTTP *

[Asserts]
status < 300

