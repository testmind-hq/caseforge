# ── GET /api/admin/users - valid request with all required fields ──
# case_id=TC-e7fb82c9
# case_name=GET /api/admin/users - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

GET {{base_url}}/api/admin/users

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.users" exists

