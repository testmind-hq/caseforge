# ── [OWASP-API8] OPTIONS /api/admin/teams/{id}/members — CORS security configuration ──
# case_id=TC-02ec7afc
# case_name=[OWASP-API8] OPTIONS /api/admin/teams/{id}/members — CORS security configuration
# step_id=step-1
# step_type=test
# technique=owasp_api_top10_spec
# priority=P0

OPTIONS {{base_url}}/api/admin/teams/{id}/members
Origin: https://evil.example.com

HTTP *

[Asserts]
header "Access-Control-Allow-Origin" != "*"

