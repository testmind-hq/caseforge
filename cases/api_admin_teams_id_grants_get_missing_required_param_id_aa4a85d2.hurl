# ── GET /api/admin/teams/{id}/grants - missing required param "id" ──
# case_id=TC-aa4a85d2
# case_name=GET /api/admin/teams/{id}/grants - missing required param "id"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

GET {{base_url}}/api/admin/teams/1/grants

HTTP 422

