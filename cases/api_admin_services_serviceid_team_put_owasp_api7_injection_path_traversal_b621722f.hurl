# ── [OWASP-API7] PUT /api/admin/services/{serviceId}/team — injection (path-traversal) ──
# case_id=TC-b621722f
# case_name=[OWASP-API7] PUT /api/admin/services/{serviceId}/team — injection (path-traversal)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PUT {{base_url}}/api/admin/services/..%2F..%2F..%2Fetc%2Fpasswd/team
```json
null
```

HTTP 400

