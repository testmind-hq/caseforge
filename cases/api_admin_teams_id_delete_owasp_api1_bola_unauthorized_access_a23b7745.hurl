# ── [OWASP-API1] DELETE /api/admin/teams/{id} — BOLA unauthorized access ──
# case_id=TC-a23b7745
# case_name=[OWASP-API1] DELETE /api/admin/teams/{id} — BOLA unauthorized access
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

DELETE {{base_url}}/api/admin/teams/{{other_resource_id}}

HTTP 403

