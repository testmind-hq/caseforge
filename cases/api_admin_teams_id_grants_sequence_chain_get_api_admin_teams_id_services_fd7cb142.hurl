# ══════════════════════════════════════════════════
# sequence chain: /api/admin/teams/{id}/grants → GET /api/admin/teams/{id}/services
# case_id=TC-fd7cb142
# case_name=sequence chain: /api/admin/teams/{id}/grants → GET /api/admin/teams/{id}/services
# case_kind=chain
# priority=P1
# ══════════════════════════════════════════════════

# ── create via POST /api/admin/teams/{id}/grants [setup] ──
# step_id=step-setup
# step_type=setup
# title=create via POST /api/admin/teams/{id}/grants

POST {{base_url}}/api/admin/teams/{id}/grants
Content-Type: application/json
```json
{
  "branches": [
    "nightly"
  ],
  "expiresAt": "2014-07-24T15:17:10Z",
  "granteeTeamId": "da38f17d-bcba-48c6-b1e9-2b8c5c84b849",
  "granteeUserId": "a204f443-d1b0-4bfc-803a-4c17ae6cc61d",
  "serviceId": "ce438324-485f-4319-9bd6-11c6d9721984"
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

