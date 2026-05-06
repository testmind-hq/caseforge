# ══════════════════════════════════════════════════
# POST /api/admin/webhooks/:id/test - idempotent: second call must be safe
# case_id=TC-ff996bd3
# case_name=POST /api/admin/webhooks/:id/test - idempotent: second call must be safe
# case_kind=chain
# priority=P2
# ══════════════════════════════════════════════════

# ── POST /api/admin/webhooks/:id/test — first call [setup] ──
# step_id=step-setup
# step_type=setup
# title=POST /api/admin/webhooks/:id/test — first call

POST {{base_url}}/api/admin/webhooks/:id/test

HTTP 200

[Asserts]
duration < 2000

# ── POST /api/admin/webhooks/:id/test — identical second call must be safe [test] ──
# step_id=step-test
# step_type=test
# title=POST /api/admin/webhooks/:id/test — identical second call must be safe
# depends_on=step-setup

POST {{base_url}}/api/admin/webhooks/:id/test

HTTP 200

[Asserts]
duration < 2000

