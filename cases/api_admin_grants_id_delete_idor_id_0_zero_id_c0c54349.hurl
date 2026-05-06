# ── DELETE /api/admin/grants/{id} - IDOR id=0 (zero_id) ──
# case_id=TC-c0c54349
# case_name=DELETE /api/admin/grants/{id} - IDOR id=0 (zero_id)
# step_id=step-main
# step_type=test
# technique=idor
# priority=P1

DELETE {{base_url}}/api/admin/grants/0

HTTP *

[Asserts]
status >= 400
status < 500

