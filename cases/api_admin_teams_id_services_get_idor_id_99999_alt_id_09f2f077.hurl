# ── GET /api/admin/teams/{id}/services - IDOR id=99999 (alt_id) ──
# case_id=TC-09f2f077
# case_name=GET /api/admin/teams/{id}/services - IDOR id=99999 (alt_id)
# step_id=step-main
# step_type=test
# technique=idor
# priority=P1

GET {{base_url}}/api/admin/teams/99999/services

HTTP *

[Asserts]
status >= 400
status < 500

