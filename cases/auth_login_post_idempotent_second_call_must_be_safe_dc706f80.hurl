# ══════════════════════════════════════════════════
# POST /auth/login - idempotent: second call must be safe
# case_id=TC-dc706f80
# case_name=POST /auth/login - idempotent: second call must be safe
# case_kind=chain
# priority=P2
# ══════════════════════════════════════════════════

# ── POST /auth/login — first call [setup] ──
# step_id=step-setup
# step_type=setup
# title=POST /auth/login — first call

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "arvidhanson@deckow.com",
  "password": "thoughtful"
}
```

HTTP 200

[Asserts]
duration < 2000

# ── POST /auth/login — identical second call must be safe [test] ──
# step_id=step-test
# step_type=test
# title=POST /auth/login — identical second call must be safe
# depends_on=step-setup

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "arvidhanson@deckow.com",
  "password": "thoughtful"
}
```

HTTP 200

[Asserts]
duration < 2000

