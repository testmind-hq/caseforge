# ── POST /api/admin/teams/{id}/grants - valid request with all required fields ──
# case_id=TC-62bccfec
# case_name=POST /api/admin/teams/{id}/grants - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

POST {{base_url}}/api/admin/teams/{id}/grants
Content-Type: application/json
```json
{
  "branches": [
    "other"
  ],
  "expiresAt": "2020-03-12T16:50:23Z",
  "granteeTeamId": "fcea5c7d-08df-4a6b-a40b-cc22936c70a6",
  "granteeUserId": "4b66d87d-2a87-436a-9cba-cbd963fe3725",
  "serviceId": "20931bd8-47ab-4a34-9161-aa0f41c54efd"
}
```

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.id" exists

