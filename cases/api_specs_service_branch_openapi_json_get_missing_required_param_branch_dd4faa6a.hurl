# ── GET /api/specs/{service}/{branch}/openapi.json - missing required param "branch" ──
# case_id=TC-dd4faa6a
# case_name=GET /api/specs/{service}/{branch}/openapi.json - missing required param "branch"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

GET {{base_url}}/api/specs/1/1/openapi.json

HTTP 422

