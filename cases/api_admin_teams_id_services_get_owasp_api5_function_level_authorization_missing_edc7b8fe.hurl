# ── [OWASP-API5] GET /api/admin/teams/{id}/services — function-level authorization missing ──
# case_id=TC-edc7b8fe
# case_name=[OWASP-API5] GET /api/admin/teams/{id}/services — function-level authorization missing
# step_id=step-1
# step_type=test
# technique=owasp_api_top10_spec
# priority=P0

GET {{base_url}}/api/admin/teams/{id}/services
Authorization: Bearer {{user_token}}

HTTP 403

