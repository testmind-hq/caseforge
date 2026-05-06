# ── [OWASP-API2] POST /auth/register — broken authentication ──
# case_id=TC-e8a47f18
# case_name=[OWASP-API2] POST /auth/register — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/auth/register

HTTP 401

