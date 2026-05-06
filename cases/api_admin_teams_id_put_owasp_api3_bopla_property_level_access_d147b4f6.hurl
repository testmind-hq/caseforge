# ── [OWASP-API3] PUT /api/admin/teams/{id} — BOPLA property-level access ──
# case_id=TC-d147b4f6
# case_name=[OWASP-API3] PUT /api/admin/teams/{id} — BOPLA property-level access
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": "Onward to better way!",
  "displayName": "moreover",
  "is_admin": true,
  "role": "admin"
}
```

HTTP 200

[Asserts]
jsonpath "$.is_admin" != true
jsonpath "$.role" != "admin"

