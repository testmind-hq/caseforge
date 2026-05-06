# ── GET /api/specs/{service}/{branch}/openapi.json - valid request with all required fields ──
# case_id=TC-e159fefe
# case_name=GET /api/specs/{service}/{branch}/openapi.json - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

GET {{base_url}}/api/specs/{service}/{branch}/openapi.json

HTTP 200

[Asserts]
duration < 2000

