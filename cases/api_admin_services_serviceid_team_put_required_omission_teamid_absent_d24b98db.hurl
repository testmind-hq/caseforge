# ── PUT /api/admin/services/{serviceId}/team - [required_omission] teamId absent ──
# case_id=TC-d24b98db
# case_name=PUT /api/admin/services/{serviceId}/team - [required_omission] teamId absent
# step_id=step-main
# step_type=test
# technique=required_omission
# priority=P2

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{}
```

HTTP *

[Asserts]
status >= 400
status < 500

