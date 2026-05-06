# ── [OWASP-API2] GET /api/admin/webhooks — broken authentication ──
# case_id=TC-ec46e5a8
# case_name=[OWASP-API2] GET /api/admin/webhooks — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/admin/webhooks

HTTP 401

