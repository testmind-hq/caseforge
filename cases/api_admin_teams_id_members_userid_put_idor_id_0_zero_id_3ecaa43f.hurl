# ── PUT /api/admin/teams/{id}/members/{userId} - IDOR id=0 (zero_id) ──
# case_id=TC-3ecaa43f
# case_name=PUT /api/admin/teams/{id}/members/{userId} - IDOR id=0 (zero_id)
# step_id=step-main
# step_type=test
# technique=idor
# priority=P1

PUT {{base_url}}/api/admin/teams/0/members/1

HTTP *

[Asserts]
status >= 400
status < 500

