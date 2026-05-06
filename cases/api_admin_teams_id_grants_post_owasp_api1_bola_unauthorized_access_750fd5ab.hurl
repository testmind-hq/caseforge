# ── [OWASP-API1] POST /api/admin/teams/{id}/grants — BOLA unauthorized access ──
# case_id=TC-750fd5ab
# case_name=[OWASP-API1] POST /api/admin/teams/{id}/grants — BOLA unauthorized access
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/teams/{{other_resource_id}}/grants

HTTP 403

