# ── GET /api/admin/teams/{id}/services - IDOR id=0 (zero_id) ──
# case_id=TC-405d2163
# case_name=GET /api/admin/teams/{id}/services - IDOR id=0 (zero_id)
# step_id=step-main
# step_type=test
# technique=idor
# priority=P1

GET {{base_url}}/api/admin/teams/0/services

HTTP *

[Asserts]
status >= 400
status < 500

