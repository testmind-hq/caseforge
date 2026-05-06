# ── [OWASP-API5] PATCH /api/admin/webhooks/:id — function-level authorization missing ──
# case_id=TC-6c16dac4
# case_name=[OWASP-API5] PATCH /api/admin/webhooks/:id — function-level authorization missing
# step_id=step-1
# step_type=test
# technique=owasp_api_top10_spec
# priority=P0

PATCH {{base_url}}/api/admin/webhooks/:id
Authorization: Bearer {{user_token}}

HTTP 403

