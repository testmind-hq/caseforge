# ── [OWASP-API1] DELETE /api/tokens/{id} — BOLA unauthorized access ──
# case_id=TC-2d207a0d
# case_name=[OWASP-API1] DELETE /api/tokens/{id} — BOLA unauthorized access
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

DELETE {{base_url}}/api/tokens/{{other_resource_id}}

HTTP 403

