# ── [OWASP-API5] GET /api/admin/users — function-level authorization missing ──
# case_id=TC-3724bb26
# case_name=[OWASP-API5] GET /api/admin/users — function-level authorization missing
# step_id=step-1
# step_type=test
# technique=owasp_api_top10_spec
# priority=P0

GET {{base_url}}/api/admin/users
Authorization: Bearer {{user_token}}

HTTP 403

