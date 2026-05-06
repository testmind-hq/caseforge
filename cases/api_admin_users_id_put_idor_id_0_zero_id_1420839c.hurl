# ── PUT /api/admin/users/{id} - IDOR id=0 (zero_id) ──
# case_id=TC-1420839c
# case_name=PUT /api/admin/users/{id} - IDOR id=0 (zero_id)
# step_id=step-main
# step_type=test
# technique=idor
# priority=P1

PUT {{base_url}}/api/admin/users/0

HTTP *

[Asserts]
status >= 400
status < 500

