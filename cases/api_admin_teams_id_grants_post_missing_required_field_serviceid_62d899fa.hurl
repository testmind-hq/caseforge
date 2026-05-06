# ── POST /api/admin/teams/{id}/grants - missing required field "serviceId" ──
# case_id=TC-62d899fa
# case_name=POST /api/admin/teams/{id}/grants - missing required field "serviceId"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

POST {{base_url}}/api/admin/teams/{id}/grants
Content-Type: application/json
```json
{
  "branches": [
    "for"
  ],
  "expiresAt": "1953-03-29T14:02:05Z",
  "granteeTeamId": "6d698330-9f66-45db-a309-61a79c0db5ba",
  "granteeUserId": "8867a80d-0d36-4338-ae27-3e2177ebe961"
}
```

HTTP 422

