# ── POST /api/admin/teams/{id}/members - IDOR id=0 (zero_id) ──
# case_id=TC-07948765
# case_name=POST /api/admin/teams/{id}/members - IDOR id=0 (zero_id)
# step_id=step-main
# step_type=test
# technique=idor
# priority=P1

POST {{base_url}}/api/admin/teams/0/members

HTTP *

[Asserts]
status >= 400
status < 500

