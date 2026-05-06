# ── [OWASP-API2] DELETE /api/admin/teams/{id} — broken authentication ──
# case_id=TC-f7305717
# case_name=[OWASP-API2] DELETE /api/admin/teams/{id} — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

DELETE {{base_url}}/api/admin/teams/{id}

HTTP 401

