# ── PUT /api/admin/services/{serviceId}/team - [semantic_annotation] nullable field "teamId" accepts null ──
# case_id=TC-f06bfa27
# case_name=PUT /api/admin/services/{serviceId}/team - [semantic_annotation] nullable field "teamId" accepts null
# step_id=step-main
# step_type=test
# technique=semantic_annotation
# priority=P1

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{
  "teamId": null
}
```

HTTP *

[Asserts]
status >= 200
status < 300

