# ══════════════════════════════════════════════════
# DELETE /api/catalog/:serviceId - idempotent: second call must be safe
# case_id=TC-84233d9e
# case_name=DELETE /api/catalog/:serviceId - idempotent: second call must be safe
# case_kind=chain
# priority=P2
# ══════════════════════════════════════════════════

# ── DELETE /api/catalog/:serviceId — first call [setup] ──
# step_id=step-setup
# step_type=setup
# title=DELETE /api/catalog/:serviceId — first call

DELETE {{base_url}}/api/catalog/:serviceId

HTTP 200

[Asserts]
duration < 2000

# ── DELETE /api/catalog/:serviceId — identical second call must be safe [test] ──
# step_id=step-test
# step_type=test
# title=DELETE /api/catalog/:serviceId — identical second call must be safe
# depends_on=step-setup

DELETE {{base_url}}/api/catalog/:serviceId

HTTP 200

[Asserts]
duration < 2000

