# ── GET /api/admin/teams/{id}/grants - valid request with all required fields ──
# case_id=TC-d5427a01
# case_name=GET /api/admin/teams/{id}/grants - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

GET {{base_url}}/api/admin/teams/{id}/grants

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.outgoing" exists
jsonpath "$.incoming" exists

