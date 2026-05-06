# ── PUT /api/admin/services/{serviceId}/team - [type_coercion] teamId wrong_type_boolean ──
# case_id=TC-5b55ebea
# case_name=PUT /api/admin/services/{serviceId}/team - [type_coercion] teamId wrong_type_boolean
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{
  "teamId": true
}
```

HTTP 422

