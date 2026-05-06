# ── POST /api/admin/teams/{id}/grants - missing required param "id" ──
# case_id=TC-aee10eee
# case_name=POST /api/admin/teams/{id}/grants - missing required param "id"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

POST {{base_url}}/api/admin/teams/1/grants

HTTP 422

