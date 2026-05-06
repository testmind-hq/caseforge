# ── [OWASP-API7] GET /api/admin/audit-logs — injection (sqli) ──
# case_id=TC-605a4d60
# case_name=[OWASP-API7] GET /api/admin/audit-logs — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/admin/audit-logs?action=%27+OR+1%3D1--
```json
null
```

HTTP 400

