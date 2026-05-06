# ── [OWASP-API2] POST /api/admin/teams/{id}/members — broken authentication ──
# case_id=TC-d1200108
# case_name=[OWASP-API2] POST /api/admin/teams/{id}/members — broken authentication
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/teams/{id}/members

HTTP 401

