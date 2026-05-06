# ── [OWASP-API2] DELETE /api/admin/users/{id} — broken authentication ──
# case_id=TC-5cc69e63
# case_name=[OWASP-API2] DELETE /api/admin/users/{id} — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

DELETE {{base_url}}/api/admin/users/{id}

HTTP 401

