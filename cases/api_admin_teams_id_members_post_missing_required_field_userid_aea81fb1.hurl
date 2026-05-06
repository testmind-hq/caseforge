# ── POST /api/admin/teams/{id}/members - missing required field "userId" ──
# case_id=TC-aea81fb1
# case_name=POST /api/admin/teams/{id}/members - missing required field "userId"
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P1

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "owner"
}
```

HTTP 422

