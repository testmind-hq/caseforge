# ── POST /api/admin/webhooks/:id/test - valid request with all required fields ──
# case_id=TC-ae0a2dc3
# case_name=POST /api/admin/webhooks/:id/test - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

POST {{base_url}}/api/admin/webhooks/:id/test

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.ok" exists

