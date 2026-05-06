# ── [OWASP-API2] POST /api/upload — broken authentication ──
# case_id=TC-4c9fd28e
# case_name=[OWASP-API2] POST /api/upload — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/upload

HTTP 401

