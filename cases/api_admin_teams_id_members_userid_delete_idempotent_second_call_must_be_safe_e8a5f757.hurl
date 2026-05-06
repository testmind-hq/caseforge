# ══════════════════════════════════════════════════
# DELETE /api/admin/teams/{id}/members/{userId} - idempotent: second call must be safe
# case_id=TC-e8a5f757
# case_name=DELETE /api/admin/teams/{id}/members/{userId} - idempotent: second call must be safe
# case_kind=chain
# priority=P2
# ══════════════════════════════════════════════════

# ── DELETE /api/admin/teams/{id}/members/{userId} — first call [setup] ──
# step_id=step-setup
# step_type=setup
# title=DELETE /api/admin/teams/{id}/members/{userId} — first call

DELETE {{base_url}}/api/admin/teams/{id}/members/{userId}

HTTP 200

[Asserts]
duration < 2000

# ── DELETE /api/admin/teams/{id}/members/{userId} — identical second call must be safe [test] ──
# step_id=step-test
# step_type=test
# title=DELETE /api/admin/teams/{id}/members/{userId} — identical second call must be safe
# depends_on=step-setup

DELETE {{base_url}}/api/admin/teams/{id}/members/{userId}

HTTP 200

[Asserts]
duration < 2000

