# ══════════════════════════════════════════════════
# PUT /api/admin/users/{id} - idempotent: second call must be safe
# case_id=TC-383d2878
# case_name=PUT /api/admin/users/{id} - idempotent: second call must be safe
# case_kind=chain
# priority=P2
# ══════════════════════════════════════════════════

# ── PUT /api/admin/users/{id} — first call [setup] ──
# step_id=step-setup
# step_type=setup
# title=PUT /api/admin/users/{id} — first call

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": false,
  "role": "team_owner"
}
```

HTTP 200

[Asserts]
duration < 2000

# ── PUT /api/admin/users/{id} — identical second call must be safe [test] ──
# step_id=step-test
# step_type=test
# title=PUT /api/admin/users/{id} — identical second call must be safe
# depends_on=step-setup

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": false,
  "role": "team_owner"
}
```

HTTP 200

[Asserts]
duration < 2000

