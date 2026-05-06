# ── [OWASP-API2] PUT /api/admin/teams/{id} — broken authentication ──
# case_id=TC-fea6c4f7
# case_name=[OWASP-API2] PUT /api/admin/teams/{id} — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PUT {{base_url}}/api/admin/teams/{id}

HTTP 401

