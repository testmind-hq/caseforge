# ══════════════════════════════════════════════════
# POST /api/tokens - idempotent: second call must be safe
# case_id=TC-85621889
# case_name=POST /api/tokens - idempotent: second call must be safe
# case_kind=chain
# priority=P2
# ══════════════════════════════════════════════════

# ── POST /api/tokens — first call [setup] ──
# step_id=step-setup
# step_type=setup
# title=POST /api/tokens — first call

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Kaya Saunders",
  "scope": "read"
}
```

HTTP 200

[Asserts]
duration < 2000

# ── POST /api/tokens — identical second call must be safe [test] ──
# step_id=step-test
# step_type=test
# title=POST /api/tokens — identical second call must be safe
# depends_on=step-setup

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Kaya Saunders",
  "scope": "read"
}
```

HTTP 200

[Asserts]
duration < 2000

