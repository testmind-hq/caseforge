# ══════════════════════════════════════════════════
# sequence chain: /api/admin/webhooks → DELETE /api/admin/users/{id}
# case_id=TC-763b85b6
# case_name=sequence chain: /api/admin/webhooks → DELETE /api/admin/users/{id}
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
    "other"
  ],
  "name": "Payton Yang",
  "providerType": "anyone",
  "teamId": "e7136d75-172b-46d0-8e7e-838fb2a645b4",
  "url": "http://www.investorarchitectures.com/viral/real-time"
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

