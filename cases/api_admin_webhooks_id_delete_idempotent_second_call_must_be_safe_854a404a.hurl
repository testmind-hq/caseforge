# ══════════════════════════════════════════════════
# DELETE /api/admin/webhooks/:id - idempotent: second call must be safe
# case_id=TC-854a404a
# case_name=DELETE /api/admin/webhooks/:id - idempotent: second call must be safe
# case_kind=chain
# priority=P2
# ══════════════════════════════════════════════════

# ── DELETE /api/admin/webhooks/:id — first call [setup] ──
# step_id=step-setup
# step_type=setup
# title=DELETE /api/admin/webhooks/:id — first call

DELETE {{base_url}}/api/admin/webhooks/:id

HTTP 204

[Asserts]
duration < 2000

# ── DELETE /api/admin/webhooks/:id — identical second call must be safe [test] ──
# step_id=step-test
# step_type=test
# title=DELETE /api/admin/webhooks/:id — identical second call must be safe
# depends_on=step-setup

DELETE {{base_url}}/api/admin/webhooks/:id

HTTP 204

[Asserts]
duration < 2000

