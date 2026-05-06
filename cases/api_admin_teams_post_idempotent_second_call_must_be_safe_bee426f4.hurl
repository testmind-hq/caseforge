# ══════════════════════════════════════════════════
# POST /api/admin/teams - idempotent: second call must be safe
# case_id=TC-bee426f4
# case_name=POST /api/admin/teams - idempotent: second call must be safe
# case_kind=chain
# priority=P2
# ══════════════════════════════════════════════════

# ── POST /api/admin/teams — first call [setup] ──
# step_id=step-setup
# step_type=setup
# title=POST /api/admin/teams — first call

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Theirs year do ready for idea.",
  "displayName": "quality",
  "name": "Lillie Hart"
}
```

HTTP 200

[Asserts]
duration < 2000

# ── POST /api/admin/teams — identical second call must be safe [test] ──
# step_id=step-test
# step_type=test
# title=POST /api/admin/teams — identical second call must be safe
# depends_on=step-setup

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Theirs year do ready for idea.",
  "displayName": "quality",
  "name": "Lillie Hart"
}
```

HTTP 200

[Asserts]
duration < 2000

