# ── GET /api/search - valid request with all required fields ──
# case_id=TC-65fdbcb4
# case_name=GET /api/search - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

GET {{base_url}}/api/search

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.results" exists

