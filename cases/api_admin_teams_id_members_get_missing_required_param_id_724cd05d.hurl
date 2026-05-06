# ── GET /api/admin/teams/{id}/members - missing required param "id" ──
# case_id=TC-724cd05d
# case_name=GET /api/admin/teams/{id}/members - missing required param "id"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

GET {{base_url}}/api/admin/teams/1/members

HTTP 422

