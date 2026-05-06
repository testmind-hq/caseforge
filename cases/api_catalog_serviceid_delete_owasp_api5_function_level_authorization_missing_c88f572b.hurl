# ── [OWASP-API5] DELETE /api/catalog/:serviceId — function-level authorization missing ──
# case_id=TC-c88f572b
# case_name=[OWASP-API5] DELETE /api/catalog/:serviceId — function-level authorization missing
# step_id=step-1
# step_type=test
# technique=owasp_api_top10_spec
# priority=P0

DELETE {{base_url}}/api/catalog/:serviceId
Authorization: Bearer {{user_token}}

HTTP 403

