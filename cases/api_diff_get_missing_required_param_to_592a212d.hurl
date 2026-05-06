# ── GET /api/diff - missing required param "to" ──
# case_id=TC-592a212d
# case_name=GET /api/diff - missing required param "to"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

GET {{base_url}}/api/diff?from=valid

HTTP 422

