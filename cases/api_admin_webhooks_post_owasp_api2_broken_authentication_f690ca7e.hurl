# ── [OWASP-API2] POST /api/admin/webhooks — broken authentication ──
# case_id=TC-f690ca7e
# case_name=[OWASP-API2] POST /api/admin/webhooks — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/webhooks

HTTP 401

