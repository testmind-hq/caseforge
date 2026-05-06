# ── GET /api/me - valid request with all required fields ──
# case_id=TC-cb06322f
# case_name=GET /api/me - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

GET {{base_url}}/api/me

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.id" exists
jsonpath "$.role" exists
jsonpath "$.teams" exists
jsonpath "$.email" exists

