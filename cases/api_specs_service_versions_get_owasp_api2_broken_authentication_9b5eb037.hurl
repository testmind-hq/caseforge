# ── [OWASP-API2] GET /api/specs/:service/versions — broken authentication ──
# case_id=TC-9b5eb037
# case_name=[OWASP-API2] GET /api/specs/:service/versions — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/specs/:service/versions

HTTP 401

