# ── GET /api/specs/:service/versions - valid request with all required fields ──
# case_id=TC-f8bdece6
# case_name=GET /api/specs/:service/versions - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

GET {{base_url}}/api/specs/:service/versions

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.versions" exists

