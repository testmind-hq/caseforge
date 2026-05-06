# ── [OWASP-API5] GET /api/admin/teams — function-level authorization missing ──
# case_id=TC-a9276ccc
# case_name=[OWASP-API5] GET /api/admin/teams — function-level authorization missing
# step_id=step-1
# step_type=test
# technique=owasp_api_top10_spec
# priority=P0

GET {{base_url}}/api/admin/teams
Authorization: Bearer {{user_token}}

HTTP 403

