# ══════════════════════════════════════════════════
# DELETE /api/admin/grants/{id} - idempotent: second call must be safe
# case_id=TC-1f6fc417
# case_name=DELETE /api/admin/grants/{id} - idempotent: second call must be safe
# case_kind=chain
# priority=P2
# ══════════════════════════════════════════════════

# ── DELETE /api/admin/grants/{id} — first call [setup] ──
# step_id=step-setup
# step_type=setup
# title=DELETE /api/admin/grants/{id} — first call

DELETE {{base_url}}/api/admin/grants/{id}

HTTP 200

[Asserts]
duration < 2000

# ── DELETE /api/admin/grants/{id} — identical second call must be safe [test] ──
# step_id=step-test
# step_type=test
# title=DELETE /api/admin/grants/{id} — identical second call must be safe
# depends_on=step-setup

DELETE {{base_url}}/api/admin/grants/{id}

HTTP 200

[Asserts]
duration < 2000

