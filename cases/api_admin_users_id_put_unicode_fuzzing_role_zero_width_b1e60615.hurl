# ── PUT /api/admin/users/{id} - [unicode_fuzzing] role zero_width ──
# case_id=TC-b1e60615
# case_name=PUT /api/admin/users/{id} - [unicode_fuzzing] role zero_width
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": false,
  "role": "​hello"
}
```

HTTP 400

