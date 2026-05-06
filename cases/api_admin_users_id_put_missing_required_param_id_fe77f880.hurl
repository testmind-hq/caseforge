# ── PUT /api/admin/users/{id} - missing required param "id" ──
# case_id=TC-fe77f880
# case_name=PUT /api/admin/users/{id} - missing required param "id"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

PUT {{base_url}}/api/admin/users/1

HTTP 422

