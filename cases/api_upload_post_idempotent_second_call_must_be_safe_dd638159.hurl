# ══════════════════════════════════════════════════
# POST /api/upload - idempotent: second call must be safe
# case_id=TC-dd638159
# case_name=POST /api/upload - idempotent: second call must be safe
# case_kind=chain
# priority=P2
# ══════════════════════════════════════════════════

# ── POST /api/upload — first call [setup] ──
# step_id=step-setup
# step_type=setup
# title=POST /api/upload — first call

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "relieved",
  "commitSha": "frequently",
  "service": "inside",
  "specContent": "east"
}
```

HTTP 200

[Asserts]
duration < 2000

# ── POST /api/upload — identical second call must be safe [test] ──
# step_id=step-test
# step_type=test
# title=POST /api/upload — identical second call must be safe
# depends_on=step-setup

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "relieved",
  "commitSha": "frequently",
  "service": "inside",
  "specContent": "east"
}
```

HTTP 200

[Asserts]
duration < 2000

