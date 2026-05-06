# ── GET /api/admin/teams/{id}/services - valid request with all required fields ──
# case_id=TC-1b69193c
# case_name=GET /api/admin/teams/{id}/services - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

GET {{base_url}}/api/admin/teams/{id}/services

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.services" exists

