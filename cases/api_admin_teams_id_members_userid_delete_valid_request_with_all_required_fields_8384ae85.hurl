# ── DELETE /api/admin/teams/{id}/members/{userId} - valid request with all required fields ──
# case_id=TC-8384ae85
# case_name=DELETE /api/admin/teams/{id}/members/{userId} - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

DELETE {{base_url}}/api/admin/teams/{id}/members/{userId}

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.ok" exists

