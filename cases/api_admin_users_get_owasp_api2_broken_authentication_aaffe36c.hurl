# ── [OWASP-API2] GET /api/admin/users — broken authentication ──
# case_id=TC-aaffe36c
# case_name=[OWASP-API2] GET /api/admin/users — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/admin/users

HTTP 401

