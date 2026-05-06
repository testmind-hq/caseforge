# ══════════════════════════════════════════════════
# sequence chain: /api/admin/webhooks → GET /api/admin/teams/{id}/grants
# case_id=TC-83289d9f
# case_name=sequence chain: /api/admin/webhooks → GET /api/admin/teams/{id}/grants
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
    "yourself"
  ],
  "name": "Janis Santos",
  "providerType": "owing",
  "teamId": "f1f952e5-15e9-4e13-9296-ebf46b9a6f04",
  "url": "http://www.corporateproductize.org/vortals"
}
```

HTTP *

[Captures]
id: jsonpath "$.id"

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

