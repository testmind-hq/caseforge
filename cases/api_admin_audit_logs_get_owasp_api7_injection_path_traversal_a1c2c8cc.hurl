# ── [OWASP-API7] GET /api/admin/audit-logs — injection (path-traversal) ──
# case_id=TC-a1c2c8cc
# case_name=[OWASP-API7] GET /api/admin/audit-logs — injection (path-traversal)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/admin/audit-logs?action=..%2F..%2F..%2Fetc%2Fpasswd
```json
null
```

HTTP 400

