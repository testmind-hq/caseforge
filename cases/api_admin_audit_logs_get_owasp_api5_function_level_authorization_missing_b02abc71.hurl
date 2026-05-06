# ── [OWASP-API5] GET /api/admin/audit-logs — function-level authorization missing ──
# case_id=TC-b02abc71
# case_name=[OWASP-API5] GET /api/admin/audit-logs — function-level authorization missing
# step_id=step-1
# step_type=test
# technique=owasp_api_top10_spec
# priority=P0

GET {{base_url}}/api/admin/audit-logs
Authorization: Bearer {{user_token}}

HTTP 403

