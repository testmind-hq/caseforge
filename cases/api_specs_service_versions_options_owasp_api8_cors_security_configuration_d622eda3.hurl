# ── [OWASP-API8] OPTIONS /api/specs/:service/versions — CORS security configuration ──
# case_id=TC-d622eda3
# case_name=[OWASP-API8] OPTIONS /api/specs/:service/versions — CORS security configuration
# step_id=step-1
# step_type=test
# technique=owasp_api_top10_spec
# priority=P0

OPTIONS {{base_url}}/api/specs/:service/versions
Origin: https://evil.example.com

HTTP *

[Asserts]
header "Access-Control-Allow-Origin" != "*"

