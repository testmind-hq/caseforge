# ── [OWASP-API3] PUT /api/admin/services/{serviceId}/team — BOPLA property-level access ──
# case_id=TC-26712b87
# case_name=[OWASP-API3] PUT /api/admin/services/{serviceId}/team — BOPLA property-level access
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{
  "is_admin": true,
  "role": "admin",
  "teamId": "da2ce66b-ccba-4bc0-b582-c8fa43a6926f"
}
```

HTTP 200

[Asserts]
jsonpath "$.is_admin" != true
jsonpath "$.role" != "admin"

