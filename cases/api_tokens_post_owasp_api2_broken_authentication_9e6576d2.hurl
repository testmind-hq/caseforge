# ── [OWASP-API2] POST /api/tokens — broken authentication ──
# case_id=TC-9e6576d2
# case_name=[OWASP-API2] POST /api/tokens — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/tokens

HTTP 401

