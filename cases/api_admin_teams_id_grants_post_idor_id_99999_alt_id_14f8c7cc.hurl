# ── POST /api/admin/teams/{id}/grants - IDOR id=99999 (alt_id) ──
# case_id=TC-14f8c7cc
# case_name=POST /api/admin/teams/{id}/grants - IDOR id=99999 (alt_id)
# step_id=step-main
# step_type=test
# technique=idor
# priority=P1

POST {{base_url}}/api/admin/teams/99999/grants

HTTP *

[Asserts]
status >= 400
status < 500

