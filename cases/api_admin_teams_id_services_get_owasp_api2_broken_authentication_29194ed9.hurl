# ── [OWASP-API2] GET /api/admin/teams/{id}/services — broken authentication ──
# case_id=TC-29194ed9
# case_name=[OWASP-API2] GET /api/admin/teams/{id}/services — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/admin/teams/{id}/services

HTTP 401

