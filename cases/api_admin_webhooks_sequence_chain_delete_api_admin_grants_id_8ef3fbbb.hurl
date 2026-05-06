# ══════════════════════════════════════════════════
# sequence chain: /api/admin/webhooks → DELETE /api/admin/grants/{id}
# case_id=TC-8ef3fbbb
# case_name=sequence chain: /api/admin/webhooks → DELETE /api/admin/grants/{id}
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
    "today"
  ],
  "name": "Abe Collier",
  "providerType": "listen",
  "teamId": "7fae1382-a4cd-4c6d-9387-4f7b3c489c4e",
  "url": "https://www.staffclicks-and-mortar.biz/monetize/monetize/initiatives"
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

