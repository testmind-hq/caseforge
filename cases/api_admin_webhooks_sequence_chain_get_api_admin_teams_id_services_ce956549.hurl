# ══════════════════════════════════════════════════
# sequence chain: /api/admin/webhooks → GET /api/admin/teams/{id}/services
# case_id=TC-ce956549
# case_name=sequence chain: /api/admin/webhooks → GET /api/admin/teams/{id}/services
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
    "does"
  ],
  "name": "Joanne Peterson",
  "providerType": "extremely",
  "teamId": "85472ea1-82f2-4e21-8559-2c86837acb46",
  "url": "http://www.nationalroi.io/integrated/integrated/target/action-items"
}
```

HTTP *

[Captures]
id: jsonpath "$.id"

[Asserts]
status < 300

# ── use via GET /api/admin/teams/{id}/services [test] ──
# step_id=step-test
# step_type=test
# title=use via GET /api/admin/teams/{id}/services
# depends_on=step-setup

GET {{base_url}}/api/admin/teams/{{id}}/services

HTTP *

[Asserts]
status < 300

