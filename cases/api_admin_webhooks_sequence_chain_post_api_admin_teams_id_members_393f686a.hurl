# ══════════════════════════════════════════════════
# sequence chain: /api/admin/webhooks → POST /api/admin/teams/{id}/members
# case_id=TC-393f686a
# case_name=sequence chain: /api/admin/webhooks → POST /api/admin/teams/{id}/members
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
    "outside"
  ],
  "name": "Marlene Jacobs",
  "providerType": "for",
  "teamId": "c8d6d6a7-3cc6-4d33-b8b1-b6c03d928bf7",
  "url": "http://www.internalbrand.info/impactful/transform/web-enabled/e-commerce"
}
```

HTTP *

[Captures]
id: jsonpath "$.id"

[Asserts]
status < 300

# ── use via POST /api/admin/teams/{id}/members [test] ──
# step_id=step-test
# step_type=test
# title=use via POST /api/admin/teams/{id}/members
# depends_on=step-setup

POST {{base_url}}/api/admin/teams/{{id}}/members
Content-Type: application/json
```json
{
  "role": "member",
  "userId": "6dc4ae45-29b7-456d-b346-b29b27cb5494"
}
```

HTTP *

[Asserts]
status < 300

