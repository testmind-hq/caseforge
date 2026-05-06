# ── [OWASP-API5] POST /api/admin/webhooks — function-level authorization missing ──
# case_id=TC-d8d5bdac
# case_name=[OWASP-API5] POST /api/admin/webhooks — function-level authorization missing
# step_id=step-1
# step_type=test
# technique=owasp_api_top10_spec
# priority=P0

POST {{base_url}}/api/admin/webhooks
Authorization: Bearer {{user_token}}

HTTP 403

