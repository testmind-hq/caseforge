# ══════════════════════════════════════════════════
# DELETE /api/tokens/{id} - idempotent: second call must be safe
# case_id=TC-ea338ec1
# case_name=DELETE /api/tokens/{id} - idempotent: second call must be safe
# case_kind=chain
# priority=P2
# ══════════════════════════════════════════════════

# ── DELETE /api/tokens/{id} — first call [setup] ──
# step_id=step-setup
# step_type=setup
# title=DELETE /api/tokens/{id} — first call

DELETE {{base_url}}/api/tokens/{id}

HTTP 200

[Asserts]
duration < 2000

# ── DELETE /api/tokens/{id} — identical second call must be safe [test] ──
# step_id=step-test
# step_type=test
# title=DELETE /api/tokens/{id} — identical second call must be safe
# depends_on=step-setup

DELETE {{base_url}}/api/tokens/{id}

HTTP 200

[Asserts]
duration < 2000

