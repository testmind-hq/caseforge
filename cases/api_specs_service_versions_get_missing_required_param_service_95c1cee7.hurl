# ── GET /api/specs/:service/versions - missing required param "service" ──
# case_id=TC-95c1cee7
# case_name=GET /api/specs/:service/versions - missing required param "service"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

GET {{base_url}}/api/specs/:service/versions?branch=valid

HTTP 422

