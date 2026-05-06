# ── PATCH /api/admin/webhooks/:id - missing required param "id" ──
# case_id=TC-8a80112e
# case_name=PATCH /api/admin/webhooks/:id - missing required param "id"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

PATCH {{base_url}}/api/admin/webhooks/:id

HTTP 422

