# ── [OWASP-API2] GET /api/admin/audit-logs — broken authentication ──
# case_id=TC-eb7a16db
# case_name=[OWASP-API2] GET /api/admin/audit-logs — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/admin/audit-logs

HTTP 401

