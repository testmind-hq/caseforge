# ── [OWASP-API5] PUT /api/admin/services/{serviceId}/team — function-level authorization missing ──
# case_id=TC-544e90d2
# case_name=[OWASP-API5] PUT /api/admin/services/{serviceId}/team — function-level authorization missing
# step_id=step-1
# step_type=test
# technique=owasp_api_top10_spec
# priority=P0

PUT {{base_url}}/api/admin/services/{serviceId}/team
Authorization: Bearer {{user_token}}

HTTP 403

