# ── [OWASP-API2] GET /api/tokens — broken authentication ──
# case_id=TC-dcecca87
# case_name=[OWASP-API2] GET /api/tokens — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/tokens

HTTP 401

