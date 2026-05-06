# ── POST /api/admin/webhooks/:id/test - missing required param "id" ──
# case_id=TC-8f3b353e
# case_name=POST /api/admin/webhooks/:id/test - missing required param "id"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

POST {{base_url}}/api/admin/webhooks/:id/test

HTTP 422

