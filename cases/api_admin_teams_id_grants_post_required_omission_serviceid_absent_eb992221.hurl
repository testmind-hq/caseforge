# ── POST /api/admin/teams/{id}/grants - [required_omission] serviceId absent ──
# case_id=TC-eb992221
# case_name=POST /api/admin/teams/{id}/grants - [required_omission] serviceId absent
# step_id=step-main
# step_type=test
# technique=required_omission
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/grants
Content-Type: application/json
```json
{
  "branches": [
    "several"
  ],
  "expiresAt": "1989-03-13T15:48:36Z",
  "granteeTeamId": "849dc625-c140-49ac-bf25-8a047cafbb78",
  "granteeUserId": "f936f656-e5c6-4646-85ad-e56be5d8778e"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

