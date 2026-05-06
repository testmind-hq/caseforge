# ── [OWASP-API2] POST /api/admin/teams/{id}/grants — broken authentication ──
# case_id=TC-a5db835c
# case_name=[OWASP-API2] POST /api/admin/teams/{id}/grants — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/teams/{id}/grants

HTTP 401

