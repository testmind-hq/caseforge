# ── [OWASP-API8] OPTIONS /api/admin/services/{serviceId}/team — CORS security configuration ──
# case_id=TC-4b672517
# case_name=[OWASP-API8] OPTIONS /api/admin/services/{serviceId}/team — CORS security configuration
# step_id=step-1
# step_type=test
# technique=owasp_api_top10_spec
# priority=P0

OPTIONS {{base_url}}/api/admin/services/{serviceId}/team
Origin: https://evil.example.com

HTTP *

[Asserts]
header "Access-Control-Allow-Origin" != "*"

