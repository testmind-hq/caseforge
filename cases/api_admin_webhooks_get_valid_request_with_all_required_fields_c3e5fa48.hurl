# ── GET /api/admin/webhooks - valid request with all required fields ──
# case_id=TC-c3e5fa48
# case_name=GET /api/admin/webhooks - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

GET {{base_url}}/api/admin/webhooks

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.webhooks" exists

