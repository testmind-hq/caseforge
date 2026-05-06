# ── GET /api/admin/audit-logs - valid request with all required fields ──
# case_id=TC-04940e9f
# case_name=GET /api/admin/audit-logs - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

GET {{base_url}}/api/admin/audit-logs

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.total" exists
jsonpath "$.logs" exists
jsonpath "$.page" exists
jsonpath "$.pageSize" exists

