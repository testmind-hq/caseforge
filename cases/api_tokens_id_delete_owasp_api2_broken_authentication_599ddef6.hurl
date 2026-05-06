# ── [OWASP-API2] DELETE /api/tokens/{id} — broken authentication ──
# case_id=TC-599ddef6
# case_name=[OWASP-API2] DELETE /api/tokens/{id} — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

DELETE {{base_url}}/api/tokens/{id}

HTTP 401

