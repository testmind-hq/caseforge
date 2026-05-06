# ── [OWASP-API7] PUT /api/admin/users/{id} — injection (path-traversal) ──
# case_id=TC-e9f5a9c9
# case_name=[OWASP-API7] PUT /api/admin/users/{id} — injection (path-traversal)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PUT {{base_url}}/api/admin/users/..%2F..%2F..%2Fetc%2Fpasswd
```json
null
```

HTTP 400

