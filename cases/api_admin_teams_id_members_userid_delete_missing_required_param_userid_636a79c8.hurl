# ── DELETE /api/admin/teams/{id}/members/{userId} - missing required param "userId" ──
# case_id=TC-636a79c8
# case_name=DELETE /api/admin/teams/{id}/members/{userId} - missing required param "userId"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

DELETE {{base_url}}/api/admin/teams/1/members/1

HTTP 422

