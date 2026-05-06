# ── [OWASP-API1] DELETE /api/admin/grants/{id} — BOLA unauthorized access ──
# case_id=TC-d8d75c69
# case_name=[OWASP-API1] DELETE /api/admin/grants/{id} — BOLA unauthorized access
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

DELETE {{base_url}}/api/admin/grants/{{other_resource_id}}

HTTP 403

