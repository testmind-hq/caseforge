# ── GET /api/tokens - valid request with all required fields ──
# case_id=TC-abcd14ab
# case_name=GET /api/tokens - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

GET {{base_url}}/api/tokens

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.tokens" exists

