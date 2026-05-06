# ══════════════════════════════════════════════════
# PUT /api/admin/services/{serviceId}/team - idempotent: second call must be safe
# case_id=TC-dc1513dd
# case_name=PUT /api/admin/services/{serviceId}/team - idempotent: second call must be safe
# case_kind=chain
# priority=P2
# ══════════════════════════════════════════════════

# ── PUT /api/admin/services/{serviceId}/team — first call [setup] ──
# step_id=step-setup
# step_type=setup
# title=PUT /api/admin/services/{serviceId}/team — first call

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{
  "teamId": "b954d030-15a4-4bc5-a0ad-c5e46e96e0a7"
}
```

HTTP 200

[Asserts]
duration < 2000

# ── PUT /api/admin/services/{serviceId}/team — identical second call must be safe [test] ──
# step_id=step-test
# step_type=test
# title=PUT /api/admin/services/{serviceId}/team — identical second call must be safe
# depends_on=step-setup

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{
  "teamId": "b954d030-15a4-4bc5-a0ad-c5e46e96e0a7"
}
```

HTTP 200

[Asserts]
duration < 2000

