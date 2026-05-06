# ── GET /api/search - missing required param "q" ──
# case_id=TC-128363b8
# case_name=GET /api/search - missing required param "q"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

GET {{base_url}}/api/search?branch=valid&service=valid

HTTP 422

