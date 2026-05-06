# ── GET /api/catalog - valid request with all required fields ──
# case_id=TC-c9b53fc1
# case_name=GET /api/catalog - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

GET {{base_url}}/api/catalog

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.services" exists

