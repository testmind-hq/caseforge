# ── DELETE /api/admin/teams/{id} - missing required param "id" ──
# case_id=TC-d700a9bc
# case_name=DELETE /api/admin/teams/{id} - missing required param "id"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

DELETE {{base_url}}/api/admin/teams/1

HTTP 422

