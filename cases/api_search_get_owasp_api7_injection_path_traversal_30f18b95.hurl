# ── [OWASP-API7] GET /api/search — injection (path-traversal) ──
# case_id=TC-30f18b95
# case_name=[OWASP-API7] GET /api/search — injection (path-traversal)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/search?q=..%2F..%2F..%2Fetc%2Fpasswd
```json
null
```

HTTP 400

