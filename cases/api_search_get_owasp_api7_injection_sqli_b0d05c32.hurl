# ── [OWASP-API7] GET /api/search — injection (sqli) ──
# case_id=TC-b0d05c32
# case_name=[OWASP-API7] GET /api/search — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/search?q=%27+OR+1%3D1--
```json
null
```

HTTP 400

