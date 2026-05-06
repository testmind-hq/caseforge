# ── DELETE /api/admin/users/{id} - missing required param "id" ──
# case_id=TC-abfeb37c
# case_name=DELETE /api/admin/users/{id} - missing required param "id"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

DELETE {{base_url}}/api/admin/users/1

HTTP 422

