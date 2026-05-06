# ── DELETE /api/admin/grants/{id} - IDOR id=99999 (alt_id) ──
# case_id=TC-b20f3be6
# case_name=DELETE /api/admin/grants/{id} - IDOR id=99999 (alt_id)
# step_id=step-main
# step_type=test
# technique=idor
# priority=P1

DELETE {{base_url}}/api/admin/grants/99999

HTTP *

[Asserts]
status >= 400
status < 500

