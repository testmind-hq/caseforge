# ── GET /api/admin/teams/{id}/services - missing required param "id" ──
# case_id=TC-bbd8e250
# case_name=GET /api/admin/teams/{id}/services - missing required param "id"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

GET {{base_url}}/api/admin/teams/1/services

HTTP 422

