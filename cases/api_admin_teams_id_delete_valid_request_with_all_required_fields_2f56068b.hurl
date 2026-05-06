# ── DELETE /api/admin/teams/{id} - valid request with all required fields ──
# case_id=TC-2f56068b
# case_name=DELETE /api/admin/teams/{id} - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

DELETE {{base_url}}/api/admin/teams/{id}

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.ok" exists

