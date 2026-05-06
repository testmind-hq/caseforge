# ── GET /api/specs/{service}/{branch}/openapi.json - missing required param "service" ──
# case_id=TC-14b52fbb
# case_name=GET /api/specs/{service}/{branch}/openapi.json - missing required param "service"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

GET {{base_url}}/api/specs/1/1/openapi.json

HTTP 422

