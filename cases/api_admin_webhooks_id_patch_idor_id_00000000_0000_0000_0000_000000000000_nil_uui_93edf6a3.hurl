# ── PATCH /api/admin/webhooks/:id - IDOR id=00000000-0000-0000-0000-000000000000 (nil_uuid) ──
# case_id=TC-93edf6a3
# case_name=PATCH /api/admin/webhooks/:id - IDOR id=00000000-0000-0000-0000-000000000000 (nil_uuid)
# step_id=step-main
# step_type=test
# technique=idor
# priority=P1

PATCH {{base_url}}/api/admin/webhooks/:id

HTTP *

[Asserts]
status >= 400
status < 500

