# ── [OWASP-API2] GET /api/me — broken authentication ──
# case_id=TC-16f4aef5
# case_name=[OWASP-API2] GET /api/me — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/me

HTTP 401

