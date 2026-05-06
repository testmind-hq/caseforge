# ── [OWASP-API5] POST /api/admin/teams/{id}/grants — function-level authorization missing ──
# case_id=TC-4c520692
# case_name=[OWASP-API5] POST /api/admin/teams/{id}/grants — function-level authorization missing
# step_id=step-1
# step_type=test
# technique=owasp_api_top10_spec
# priority=P0

POST {{base_url}}/api/admin/teams/{id}/grants
Authorization: Bearer {{user_token}}

HTTP 403

