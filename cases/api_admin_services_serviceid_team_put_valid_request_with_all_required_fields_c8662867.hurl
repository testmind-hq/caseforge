# ── PUT /api/admin/services/{serviceId}/team - valid request with all required fields ──
# case_id=TC-c8662867
# case_name=PUT /api/admin/services/{serviceId}/team - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{
  "teamId": "8439a10e-558d-4569-b260-f0f36a116d83"
}
```

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.ok" exists

