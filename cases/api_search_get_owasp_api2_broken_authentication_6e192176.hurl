# ── [OWASP-API2] GET /api/search — broken authentication ──
# case_id=TC-6e192176
# case_name=[OWASP-API2] GET /api/search — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/search

HTTP 401

