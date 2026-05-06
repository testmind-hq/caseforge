# ══════════════════════════════════════════════════
# POST /api/admin/webhooks - idempotent: second call must be safe
# case_id=TC-06e188f6
# case_name=POST /api/admin/webhooks - idempotent: second call must be safe
# case_kind=chain
# priority=P2
# ══════════════════════════════════════════════════

# ── POST /api/admin/webhooks — first call [setup] ──
# step_id=step-setup
# step_type=setup
# title=POST /api/admin/webhooks — first call

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": [
    "now"
  ],
  "name": "Anya Wright",
  "providerType": "yesterday",
  "teamId": "cd7a7947-5e97-4e0c-bd41-40373e8f332b",
  "url": "http://www.primaryaction-items.org/enhance/deploy/interfaces"
}
```

HTTP 200

[Asserts]
duration < 2000

# ── POST /api/admin/webhooks — identical second call must be safe [test] ──
# step_id=step-test
# step_type=test
# title=POST /api/admin/webhooks — identical second call must be safe
# depends_on=step-setup

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": [
    "now"
  ],
  "name": "Anya Wright",
  "providerType": "yesterday",
  "teamId": "cd7a7947-5e97-4e0c-bd41-40373e8f332b",
  "url": "http://www.primaryaction-items.org/enhance/deploy/interfaces"
}
```

HTTP 200

[Asserts]
duration < 2000

