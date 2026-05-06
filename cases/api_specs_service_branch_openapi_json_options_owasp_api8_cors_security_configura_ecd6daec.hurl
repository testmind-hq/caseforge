# ── [OWASP-API8] OPTIONS /api/specs/{service}/{branch}/openapi.json — CORS security configuration ──
# case_id=TC-ecd6daec
# case_name=[OWASP-API8] OPTIONS /api/specs/{service}/{branch}/openapi.json — CORS security configuration
# step_id=step-1
# step_type=test
# technique=owasp_api_top10_spec
# priority=P0

OPTIONS {{base_url}}/api/specs/{service}/{branch}/openapi.json
Origin: https://evil.example.com

HTTP *

[Asserts]
header "Access-Control-Allow-Origin" != "*"

