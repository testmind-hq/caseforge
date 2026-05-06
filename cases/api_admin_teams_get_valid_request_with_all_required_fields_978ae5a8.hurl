# ── GET /api/admin/teams - valid request with all required fields ──
# case_id=TC-978ae5a8
# case_name=GET /api/admin/teams - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

GET {{base_url}}/api/admin/teams

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.teams" exists

