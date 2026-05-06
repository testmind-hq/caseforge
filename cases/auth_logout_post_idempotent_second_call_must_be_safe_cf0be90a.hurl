# ══════════════════════════════════════════════════
# POST /auth/logout - idempotent: second call must be safe
# case_id=TC-cf0be90a
# case_name=POST /auth/logout - idempotent: second call must be safe
# case_kind=chain
# priority=P2
# ══════════════════════════════════════════════════

# ── POST /auth/logout — first call [setup] ──
# step_id=step-setup
# step_type=setup
# title=POST /auth/logout — first call

POST {{base_url}}/auth/logout

HTTP 200

[Asserts]
duration < 2000

# ── POST /auth/logout — identical second call must be safe [test] ──
# step_id=step-test
# step_type=test
# title=POST /auth/logout — identical second call must be safe
# depends_on=step-setup

POST {{base_url}}/auth/logout

HTTP 200

[Asserts]
duration < 2000

