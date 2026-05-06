# ── GET /api/admin/teams/{id}/members - IDOR id=0 (zero_id) ──
# case_id=TC-8d769a8b
# case_name=GET /api/admin/teams/{id}/members - IDOR id=0 (zero_id)
# step_id=step-main
# step_type=test
# technique=idor
# priority=P1

GET {{base_url}}/api/admin/teams/0/members

HTTP *

[Asserts]
status >= 400
status < 500

