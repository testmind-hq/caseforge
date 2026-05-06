# ══════════════════════════════════════════════════
# POST /api/admin/teams/{id}/grants - idempotent: second call must be safe
# case_id=TC-810053e8
# case_name=POST /api/admin/teams/{id}/grants - idempotent: second call must be safe
# case_kind=chain
# priority=P2
# ══════════════════════════════════════════════════

# ── POST /api/admin/teams/{id}/grants — first call [setup] ──
# step_id=step-setup
# step_type=setup
# title=POST /api/admin/teams/{id}/grants — first call

POST {{base_url}}/api/admin/teams/{id}/grants
Content-Type: application/json
```json
{
  "branches": [
    "mercy"
  ],
  "expiresAt": "1999-12-17T23:28:47Z",
  "granteeTeamId": "65e38a66-d932-4217-b7b6-b9d191c81aaf",
  "granteeUserId": "41f62f9a-dcd8-4b25-86af-1c3d9ec30857",
  "serviceId": "4926c858-e08e-4a3f-bf7b-0bb8e4309181"
}
```

HTTP 200

[Asserts]
duration < 2000

# ── POST /api/admin/teams/{id}/grants — identical second call must be safe [test] ──
# step_id=step-test
# step_type=test
# title=POST /api/admin/teams/{id}/grants — identical second call must be safe
# depends_on=step-setup

POST {{base_url}}/api/admin/teams/{id}/grants
Content-Type: application/json
```json
{
  "branches": [
    "mercy"
  ],
  "expiresAt": "1999-12-17T23:28:47Z",
  "granteeTeamId": "65e38a66-d932-4217-b7b6-b9d191c81aaf",
  "granteeUserId": "41f62f9a-dcd8-4b25-86af-1c3d9ec30857",
  "serviceId": "4926c858-e08e-4a3f-bf7b-0bb8e4309181"
}
```

HTTP 200

[Asserts]
duration < 2000

