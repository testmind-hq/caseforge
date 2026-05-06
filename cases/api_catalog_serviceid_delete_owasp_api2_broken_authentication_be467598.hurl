# ── [OWASP-API2] DELETE /api/catalog/:serviceId — broken authentication ──
# case_id=TC-be467598
# case_name=[OWASP-API2] DELETE /api/catalog/:serviceId — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

DELETE {{base_url}}/api/catalog/:serviceId

HTTP 401

