# ══════════════════════════════════════════════════
# sequence chain: /api/admin/webhooks → PUT /api/admin/services/{serviceId}/team
# case_id=TC-256209eb
# case_name=sequence chain: /api/admin/webhooks → PUT /api/admin/services/{serviceId}/team
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
    "throughout"
  ],
  "name": "Pablo Hoffman",
  "providerType": "barely",
  "teamId": "cc3b8d87-6c30-464d-a451-ec70a317a56a",
  "url": "http://www.futuresynergize.org/evolve"
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
  "teamId": "fbaecfc9-d46e-4518-8fc8-3534e881b114"
}
```

HTTP *

[Asserts]
status < 300

