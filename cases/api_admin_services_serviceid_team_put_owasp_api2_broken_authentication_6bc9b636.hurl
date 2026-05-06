# ── [OWASP-API2] PUT /api/admin/services/{serviceId}/team — broken authentication ──
# case_id=TC-6bc9b636
# case_name=[OWASP-API2] PUT /api/admin/services/{serviceId}/team — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PUT {{base_url}}/api/admin/services/{serviceId}/team

HTTP 401

