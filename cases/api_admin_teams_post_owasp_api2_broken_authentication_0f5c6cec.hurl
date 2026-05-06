# ── [OWASP-API2] POST /api/admin/teams — broken authentication ──
# case_id=TC-0f5c6cec
# case_name=[OWASP-API2] POST /api/admin/teams — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/teams

HTTP 401

