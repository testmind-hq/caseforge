# ── [OWASP-API5] POST /api/admin/teams — function-level authorization missing ──
# case_id=TC-2df9f5ad
# case_name=[OWASP-API5] POST /api/admin/teams — function-level authorization missing
# step_id=step-1
# step_type=test
# technique=owasp_api_top10_spec
# priority=P0

POST {{base_url}}/api/admin/teams
Authorization: Bearer {{user_token}}

HTTP 403

