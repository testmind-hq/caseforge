# ── POST /api/admin/teams/{id}/members - wrong content-type (text/plain) ──
# case_id=TC-0f904569
# case_name=POST /api/admin/teams/{id}/members - wrong content-type (text/plain)
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: text/plain
```json
{
  "role": "owner",
  "userId": "b6f51cc4-2389-42c5-a864-35545c08cda9"
}
```

HTTP 415

