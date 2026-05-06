# ══════════════════════════════════════════════════
# sequence chain: /api/admin/webhooks → GET /api/admin/teams/{id}/members
# case_id=TC-969a9fae
# case_name=sequence chain: /api/admin/webhooks → GET /api/admin/teams/{id}/members
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
    "another"
  ],
  "name": "Roxanne Barber",
  "providerType": "well",
  "teamId": "360fddbd-2bf8-4533-b759-353946ddb3bb",
  "url": "https://www.corporateimplement.net/recontextualize/extensible/leading-edge"
}
```

HTTP *

[Captures]
id: jsonpath "$.id"

[Asserts]
status < 300

# ── use via GET /api/admin/teams/{id}/members [test] ──
# step_id=step-test
# step_type=test
# title=use via GET /api/admin/teams/{id}/members
# depends_on=step-setup

GET {{base_url}}/api/admin/teams/{{id}}/members

HTTP *

[Asserts]
status < 300

