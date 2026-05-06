# ── [OWASP-API1] PUT /api/admin/services/{serviceId}/team — BOLA unauthorized access ──
# case_id=TC-b7125bf5
# case_name=[OWASP-API1] PUT /api/admin/services/{serviceId}/team — BOLA unauthorized access
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PUT {{base_url}}/api/admin/services/{{other_resource_id}}/team

HTTP 403

