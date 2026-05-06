# ── POST /api/admin/teams/{id}/members - missing required param "id" ──
# case_id=TC-e44fc900
# case_name=POST /api/admin/teams/{id}/members - missing required param "id"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

POST {{base_url}}/api/admin/teams/1/members

HTTP 422

