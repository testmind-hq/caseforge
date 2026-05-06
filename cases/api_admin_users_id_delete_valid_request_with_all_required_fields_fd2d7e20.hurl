# ── DELETE /api/admin/users/{id} - valid request with all required fields ──
# case_id=TC-fd2d7e20
# case_name=DELETE /api/admin/users/{id} - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

DELETE {{base_url}}/api/admin/users/{id}

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.ok" exists

