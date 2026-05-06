# ── [OWASP-API2] GET /api/admin/teams — broken authentication ──
# case_id=TC-1e347647
# case_name=[OWASP-API2] GET /api/admin/teams — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/admin/teams

HTTP 401

