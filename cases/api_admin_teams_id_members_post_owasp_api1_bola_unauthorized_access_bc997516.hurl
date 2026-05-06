# ── [OWASP-API1] POST /api/admin/teams/{id}/members — BOLA unauthorized access ──
# case_id=TC-bc997516
# case_name=[OWASP-API1] POST /api/admin/teams/{id}/members — BOLA unauthorized access
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/teams/{{other_resource_id}}/members

HTTP 403

