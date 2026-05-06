# ── PUT /api/admin/teams/{id}/members/{userId} - wrong content-type (text/plain) ──
# case_id=TC-55f30d0f
# case_name=PUT /api/admin/teams/{id}/members/{userId} - wrong content-type (text/plain)
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: text/plain
```json
{
  "role": "member"
}
```

HTTP 415

