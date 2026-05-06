# ── PUT /api/admin/teams/{id} - missing required param "id" ──
# case_id=TC-09825850
# case_name=PUT /api/admin/teams/{id} - missing required param "id"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

PUT {{base_url}}/api/admin/teams/1

HTTP 422

