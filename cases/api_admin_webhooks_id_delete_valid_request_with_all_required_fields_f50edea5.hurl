# ── DELETE /api/admin/webhooks/:id - valid request with all required fields ──
# case_id=TC-f50edea5
# case_name=DELETE /api/admin/webhooks/:id - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

DELETE {{base_url}}/api/admin/webhooks/:id

HTTP 204

[Asserts]
duration < 2000

