# ══════════════════════════════════════════════════
# POST /api/admin/teams/{id}/members - idempotent: second call must be safe
# case_id=TC-fce8d8db
# case_name=POST /api/admin/teams/{id}/members - idempotent: second call must be safe
# case_kind=chain
# priority=P2
# ══════════════════════════════════════════════════

# ── POST /api/admin/teams/{id}/members — first call [setup] ──
# step_id=step-setup
# step_type=setup
# title=POST /api/admin/teams/{id}/members — first call

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "member",
  "userId": "f78fd0f2-6376-4a2b-8124-8006f5d96d4a"
}
```

HTTP 200

[Asserts]
duration < 2000

# ── POST /api/admin/teams/{id}/members — identical second call must be safe [test] ──
# step_id=step-test
# step_type=test
# title=POST /api/admin/teams/{id}/members — identical second call must be safe
# depends_on=step-setup

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "member",
  "userId": "f78fd0f2-6376-4a2b-8124-8006f5d96d4a"
}
```

HTTP 200

[Asserts]
duration < 2000

