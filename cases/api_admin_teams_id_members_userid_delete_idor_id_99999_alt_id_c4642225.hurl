# ── DELETE /api/admin/teams/{id}/members/{userId} - IDOR id=99999 (alt_id) ──
# case_id=TC-c4642225
# case_name=DELETE /api/admin/teams/{id}/members/{userId} - IDOR id=99999 (alt_id)
# step_id=step-main
# step_type=test
# technique=idor
# priority=P1

DELETE {{base_url}}/api/admin/teams/99999/members/1

HTTP *

[Asserts]
status >= 400
status < 500

