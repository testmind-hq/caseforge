# ── PUT /api/admin/teams/{id}/members/{userId} - missing required param "userId" ──
# case_id=TC-a0b457a0
# case_name=PUT /api/admin/teams/{id}/members/{userId} - missing required param "userId"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

PUT {{base_url}}/api/admin/teams/1/members/1

HTTP 422

