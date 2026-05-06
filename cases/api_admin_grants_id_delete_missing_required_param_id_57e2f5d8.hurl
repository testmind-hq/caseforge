# ── DELETE /api/admin/grants/{id} - missing required param "id" ──
# case_id=TC-57e2f5d8
# case_name=DELETE /api/admin/grants/{id} - missing required param "id"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

DELETE {{base_url}}/api/admin/grants/1

HTTP 422

