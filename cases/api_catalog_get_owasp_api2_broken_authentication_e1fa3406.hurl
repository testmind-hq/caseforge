# ── [OWASP-API2] GET /api/catalog — broken authentication ──
# case_id=TC-e1fa3406
# case_name=[OWASP-API2] GET /api/catalog — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/catalog

HTTP 401

