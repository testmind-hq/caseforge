# ── [OWASP-API2] DELETE /api/admin/webhooks/:id — broken authentication ──
# case_id=TC-23cf0c86
# case_name=[OWASP-API2] DELETE /api/admin/webhooks/:id — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

DELETE {{base_url}}/api/admin/webhooks/:id

HTTP 401

