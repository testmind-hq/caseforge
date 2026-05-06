# ── [OWASP-API5] POST /api/admin/webhooks/:id/test — function-level authorization missing ──
# case_id=TC-908d0d93
# case_name=[OWASP-API5] POST /api/admin/webhooks/:id/test — function-level authorization missing
# step_id=step-1
# step_type=test
# technique=owasp_api_top10_spec
# priority=P0

POST {{base_url}}/api/admin/webhooks/:id/test
Authorization: Bearer {{user_token}}

HTTP 403

