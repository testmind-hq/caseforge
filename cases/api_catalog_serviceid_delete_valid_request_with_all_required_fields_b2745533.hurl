# ── DELETE /api/catalog/:serviceId - valid request with all required fields ──
# case_id=TC-b2745533
# case_name=DELETE /api/catalog/:serviceId - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

DELETE {{base_url}}/api/catalog/:serviceId

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.ok" exists

