# ── [OWASP-API2] GET /api/diff — broken authentication ──
# case_id=TC-f6e6d81e
# case_name=[OWASP-API2] GET /api/diff — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/diff

HTTP 401

