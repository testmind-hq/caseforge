# ── GET /api/admin/teams/{id}/members - IDOR id=99999 (alt_id) ──
# case_id=TC-4af55f13
# case_name=GET /api/admin/teams/{id}/members - IDOR id=99999 (alt_id)
# step_id=step-main
# step_type=test
# technique=idor
# priority=P1

GET {{base_url}}/api/admin/teams/99999/members

HTTP *

[Asserts]
status >= 400
status < 500

