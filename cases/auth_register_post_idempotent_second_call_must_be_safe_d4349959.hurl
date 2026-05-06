# ══════════════════════════════════════════════════
# POST /auth/register - idempotent: second call must be safe
# case_id=TC-d4349959
# case_name=POST /auth/register - idempotent: second call must be safe
# case_kind=chain
# priority=P2
# ══════════════════════════════════════════════════

# ── POST /auth/register — first call [setup] ──
# step_id=step-setup
# step_type=setup
# title=POST /auth/register — first call

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "selenagarza@ross.name",
  "password": "break"
}
```

HTTP 200

[Asserts]
duration < 2000

# ── POST /auth/register — identical second call must be safe [test] ──
# step_id=step-test
# step_type=test
# title=POST /auth/register — identical second call must be safe
# depends_on=step-setup

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "selenagarza@ross.name",
  "password": "break"
}
```

HTTP 200

[Asserts]
duration < 2000

