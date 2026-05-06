# ── PUT /api/admin/services/{serviceId}/team - wrong content-type (text/plain) ──
# case_id=TC-16d39238
# case_name=PUT /api/admin/services/{serviceId}/team - wrong content-type (text/plain)
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: text/plain
```json
{
  "teamId": "bc1c5a2f-34be-4a46-bc1a-a3abfe061eb1"
}
```

HTTP 415

