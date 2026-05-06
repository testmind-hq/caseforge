# ── POST /api/admin/teams/{id}/grants - missing required field "serviceId" ──
# case_id=TC-33636c2c
# case_name=POST /api/admin/teams/{id}/grants - missing required field "serviceId"
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P1

POST {{base_url}}/api/admin/teams/{id}/grants
Content-Type: application/json
```json
{
  "branches": [
    "lastly"
  ],
  "expiresAt": "2010-02-21T09:42:07Z",
  "granteeTeamId": "54d614e8-78c4-4be4-8d58-6262bc0ed601",
  "granteeUserId": "ebe6434a-7451-43df-a2a8-4ff4abc09840"
}
```

HTTP 422

