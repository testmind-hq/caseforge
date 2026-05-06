# ── [OWASP-API5] PUT /api/admin/teams/{id} — function-level authorization missing ──
# case_id=TC-06b71a7c
# case_name=[OWASP-API5] PUT /api/admin/teams/{id} — function-level authorization missing
# step_id=step-1
# step_type=test
# technique=owasp_api_top10_spec
# priority=P0

PUT {{base_url}}/api/admin/teams/{id}
Authorization: Bearer {{user_token}}

HTTP 403

