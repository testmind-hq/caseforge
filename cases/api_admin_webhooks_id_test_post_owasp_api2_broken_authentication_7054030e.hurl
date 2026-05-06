# ── [OWASP-API2] POST /api/admin/webhooks/:id/test — broken authentication ──
# case_id=TC-7054030e
# case_name=[OWASP-API2] POST /api/admin/webhooks/:id/test — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/webhooks/:id/test

HTTP 401

