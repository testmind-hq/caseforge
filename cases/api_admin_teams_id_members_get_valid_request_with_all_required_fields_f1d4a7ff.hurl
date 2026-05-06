# ── GET /api/admin/teams/{id}/members - valid request with all required fields ──
# case_id=TC-f1d4a7ff
# case_name=GET /api/admin/teams/{id}/members - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

GET {{base_url}}/api/admin/teams/{id}/members

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.members" exists

