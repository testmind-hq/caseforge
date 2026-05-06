# ── POST /auth/logout - valid request with all required fields ──
# case_id=TC-a517ccf9
# case_name=POST /auth/logout - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

POST {{base_url}}/auth/logout

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.ok" exists

