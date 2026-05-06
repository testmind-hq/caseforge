# ── GET /api/diff - valid request with all required fields ──
# case_id=TC-f98b2b82
# case_name=GET /api/diff - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

GET {{base_url}}/api/diff

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.added" exists
jsonpath "$.modified" exists
jsonpath "$.removed" exists

