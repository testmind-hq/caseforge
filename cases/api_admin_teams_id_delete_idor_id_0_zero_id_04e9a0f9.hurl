# ── DELETE /api/admin/teams/{id} - IDOR id=0 (zero_id) ──
# case_id=TC-04e9a0f9
# case_name=DELETE /api/admin/teams/{id} - IDOR id=0 (zero_id)
# step_id=step-main
# step_type=test
# technique=idor
# priority=P1

DELETE {{base_url}}/api/admin/teams/0

HTTP *

[Asserts]
status >= 400
status < 500

