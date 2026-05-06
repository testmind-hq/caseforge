# ── PUT /api/admin/services/{serviceId}/team - [type_coercion] teamId wrong_type_integer ──
# case_id=TC-87eccc15
# case_name=PUT /api/admin/services/{serviceId}/team - [type_coercion] teamId wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{
  "teamId": 123
}
```

HTTP 422

