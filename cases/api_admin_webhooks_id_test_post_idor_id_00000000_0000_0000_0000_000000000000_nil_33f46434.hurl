# ── POST /api/admin/webhooks/:id/test - IDOR id=00000000-0000-0000-0000-000000000000 (nil_uuid) ──
# case_id=TC-33f46434
# case_name=POST /api/admin/webhooks/:id/test - IDOR id=00000000-0000-0000-0000-000000000000 (nil_uuid)
# step_id=step-main
# step_type=test
# technique=idor
# priority=P1

POST {{base_url}}/api/admin/webhooks/:id/test

HTTP *

[Asserts]
status >= 400
status < 500

