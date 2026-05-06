# ── [OWASP-API2] PUT /api/admin/teams/{id}/members/{userId} — broken authentication ──
# case_id=TC-19b34217
# case_name=[OWASP-API2] PUT /api/admin/teams/{id}/members/{userId} — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}

HTTP 401

