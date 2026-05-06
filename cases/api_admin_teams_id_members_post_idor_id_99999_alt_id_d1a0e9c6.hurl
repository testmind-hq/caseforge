# ── POST /api/admin/teams/{id}/members - IDOR id=99999 (alt_id) ──
# case_id=TC-d1a0e9c6
# case_name=POST /api/admin/teams/{id}/members - IDOR id=99999 (alt_id)
# step_id=step-main
# step_type=test
# technique=idor
# priority=P1

POST {{base_url}}/api/admin/teams/99999/members

HTTP *

[Asserts]
status >= 400
status < 500

