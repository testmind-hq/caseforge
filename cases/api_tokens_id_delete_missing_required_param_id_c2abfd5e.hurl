# ── DELETE /api/tokens/{id} - missing required param "id" ──
# case_id=TC-c2abfd5e
# case_name=DELETE /api/tokens/{id} - missing required param "id"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

DELETE {{base_url}}/api/tokens/1

HTTP 422

