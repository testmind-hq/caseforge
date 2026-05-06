# ── DELETE /api/admin/webhooks/:id - missing required param "id" ──
# case_id=TC-25ba00ae
# case_name=DELETE /api/admin/webhooks/:id - missing required param "id"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

DELETE {{base_url}}/api/admin/webhooks/:id

HTTP 422

