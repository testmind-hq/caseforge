# ── [OWASP-API2] GET /api/admin/teams/{id}/grants — broken authentication ──
# case_id=TC-2dae98a0
# case_name=[OWASP-API2] GET /api/admin/teams/{id}/grants — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/admin/teams/{id}/grants

HTTP 401

