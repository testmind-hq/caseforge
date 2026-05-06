# ── DELETE /api/tokens/{id} - IDOR id=99999 (alt_id) ──
# case_id=TC-502920f7
# case_name=DELETE /api/tokens/{id} - IDOR id=99999 (alt_id)
# step_id=step-main
# step_type=test
# technique=idor
# priority=P1

DELETE {{base_url}}/api/tokens/99999

HTTP *

[Asserts]
status >= 400
status < 500

