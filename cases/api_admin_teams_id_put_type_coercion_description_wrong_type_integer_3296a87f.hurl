# ── PUT /api/admin/teams/{id} - [type_coercion] description wrong_type_integer ──
# case_id=TC-3296a87f
# case_name=PUT /api/admin/teams/{id} - [type_coercion] description wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": 123,
  "displayName": "addition"
}
```

HTTP 422

