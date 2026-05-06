# ── DELETE /api/admin/webhooks/:id - IDOR id=00000000-0000-0000-0000-000000000001 (alt_uuid) ──
# case_id=TC-101b67d9
# case_name=DELETE /api/admin/webhooks/:id - IDOR id=00000000-0000-0000-0000-000000000001 (alt_uuid)
# step_id=step-main
# step_type=test
# technique=idor
# priority=P1

DELETE {{base_url}}/api/admin/webhooks/:id

HTTP *

[Asserts]
status >= 400
status < 500

