# ══════════════════════════════════════════════════
# sequence chain: /api/admin/webhooks → PUT /api/admin/users/{id}
# case_id=TC-88a6983e
# case_name=sequence chain: /api/admin/webhooks → PUT /api/admin/users/{id}
# case_kind=chain
# priority=P1
# ══════════════════════════════════════════════════

# ── create via POST /api/admin/webhooks [setup] ──
# step_id=step-setup
# step_type=setup
# title=create via POST /api/admin/webhooks

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": [
    "only"
  ],
  "name": "Dawson Matthews",
  "providerType": "that",
  "teamId": "7c2b8aba-98b4-477e-b7fe-f53f6306f514",
  "url": "http://www.financecultivate.com/envisioneer/enable/synergies/strategize"
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
  "isActive": false,
  "role": "super_admin"
}
```

HTTP *

[Asserts]
status < 300

