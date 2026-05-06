# ── [OWASP-API2] GET /api/specs/{service}/{branch}/openapi.json — broken authentication ──
# case_id=TC-5b840153
# case_name=[OWASP-API2] GET /api/specs/{service}/{branch}/openapi.json — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/specs/{service}/{branch}/openapi.json

HTTP 401

