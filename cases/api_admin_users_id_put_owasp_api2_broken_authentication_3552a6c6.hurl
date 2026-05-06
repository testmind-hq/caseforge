# ── [OWASP-API2] PUT /api/admin/users/{id} — broken authentication ──
# case_id=TC-3552a6c6
# case_name=[OWASP-API2] PUT /api/admin/users/{id} — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PUT {{base_url}}/api/admin/users/{id}

HTTP 401

