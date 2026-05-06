# ── [OWASP-API5] GET /api/admin/webhooks — function-level authorization missing ──
# case_id=TC-a2ef426c
# case_name=[OWASP-API5] GET /api/admin/webhooks — function-level authorization missing
# step_id=step-1
# step_type=test
# technique=owasp_api_top10_spec
# priority=P0

GET {{base_url}}/api/admin/webhooks
Authorization: Bearer {{user_token}}

HTTP 403

