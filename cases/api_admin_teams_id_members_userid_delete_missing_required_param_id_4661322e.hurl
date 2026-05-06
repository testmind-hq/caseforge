# ── DELETE /api/admin/teams/{id}/members/{userId} - missing required param "id" ──
# case_id=TC-4661322e
# case_name=DELETE /api/admin/teams/{id}/members/{userId} - missing required param "id"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

DELETE {{base_url}}/api/admin/teams/1/members/1

HTTP 422

