# ── GET /api/diff - missing required param "from" ──
# case_id=TC-436315da
# case_name=GET /api/diff - missing required param "from"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

GET {{base_url}}/api/diff?to=valid

HTTP 422

