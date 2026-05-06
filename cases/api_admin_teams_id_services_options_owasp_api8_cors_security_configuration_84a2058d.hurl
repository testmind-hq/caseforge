# ── [OWASP-API8] OPTIONS /api/admin/teams/{id}/services — CORS security configuration ──
# case_id=TC-84a2058d
# case_name=[OWASP-API8] OPTIONS /api/admin/teams/{id}/services — CORS security configuration
# step_id=step-1
# step_type=test
# technique=owasp_api_top10_spec
# priority=P0

OPTIONS {{base_url}}/api/admin/teams/{id}/services
Origin: https://evil.example.com

HTTP *

[Asserts]
header "Access-Control-Allow-Origin" != "*"

