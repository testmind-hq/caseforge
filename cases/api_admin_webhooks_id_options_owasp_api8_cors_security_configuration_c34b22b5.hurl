# ── [OWASP-API8] OPTIONS /api/admin/webhooks/:id — CORS security configuration ──
# case_id=TC-c34b22b5
# case_name=[OWASP-API8] OPTIONS /api/admin/webhooks/:id — CORS security configuration
# step_id=step-1
# step_type=test
# technique=owasp_api_top10_spec
# priority=P0

OPTIONS {{base_url}}/api/admin/webhooks/:id
Origin: https://evil.example.com

HTTP *

[Asserts]
header "Access-Control-Allow-Origin" != "*"

