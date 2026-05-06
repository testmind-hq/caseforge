# ══════════════════════════════════════════════════
# sequence chain: /api/admin/webhooks → POST /api/admin/teams/{id}/grants
# case_id=TC-02ba968a
# case_name=sequence chain: /api/admin/webhooks → POST /api/admin/teams/{id}/grants
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
    "him"
  ],
  "name": "Cayla Rosenbaum",
  "providerType": "ours",
  "teamId": "ccd3929e-a106-4df3-8d31-66697e80dbe3",
  "url": "https://www.seniore-enable.name/synergies/end-to-end/integrate/e-tailers"
}
```

HTTP *

[Captures]
id: jsonpath "$.id"

[Asserts]
status < 300

# ── use via POST /api/admin/teams/{id}/grants [test] ──
# step_id=step-test
# step_type=test
# title=use via POST /api/admin/teams/{id}/grants
# depends_on=step-setup

POST {{base_url}}/api/admin/teams/{{id}}/grants
Content-Type: application/json
```json
{
  "branches": [
    "i.e."
  ],
  "expiresAt": "2011-10-23T02:54:47Z",
  "granteeTeamId": "d189b00e-5719-4cc5-b97a-a00f62029da1",
  "granteeUserId": "77c00823-081e-4450-9ea4-1bd04aabfdee",
  "serviceId": "433f7b49-b2b9-485d-a48e-d48715ed6be5"
}
```

HTTP *

[Asserts]
status < 300

