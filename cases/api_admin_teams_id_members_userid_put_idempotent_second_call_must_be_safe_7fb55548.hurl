# ══════════════════════════════════════════════════
# PUT /api/admin/teams/{id}/members/{userId} - idempotent: second call must be safe
# case_id=TC-7fb55548
# case_name=PUT /api/admin/teams/{id}/members/{userId} - idempotent: second call must be safe
# case_kind=chain
# priority=P2
# ══════════════════════════════════════════════════

# ── PUT /api/admin/teams/{id}/members/{userId} — first call [setup] ──
# step_id=step-setup
# step_type=setup
# title=PUT /api/admin/teams/{id}/members/{userId} — first call

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "role": "owner"
}
```

HTTP 200

[Asserts]
duration < 2000

# ── PUT /api/admin/teams/{id}/members/{userId} — identical second call must be safe [test] ──
# step_id=step-test
# step_type=test
# title=PUT /api/admin/teams/{id}/members/{userId} — identical second call must be safe
# depends_on=step-setup

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "role": "owner"
}
```

HTTP 200

[Asserts]
duration < 2000

