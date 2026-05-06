# ── DELETE /api/tokens/{id} - IDOR id=0 (zero_id) ──
# case_id=TC-d0e0481e
# case_name=DELETE /api/tokens/{id} - IDOR id=0 (zero_id)
# step_id=step-main
# step_type=test
# technique=idor
# priority=P1

DELETE {{base_url}}/api/tokens/0

HTTP *

[Asserts]
status >= 400
status < 500

