# ── GET /api/specs/:service/versions - missing required param "branch" ──
# case_id=TC-e71dd727
# case_name=GET /api/specs/:service/versions - missing required param "branch"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

GET {{base_url}}/api/specs/:service/versions

HTTP 422

