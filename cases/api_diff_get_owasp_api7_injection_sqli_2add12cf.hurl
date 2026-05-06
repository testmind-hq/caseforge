# ── [OWASP-API7] GET /api/diff — injection (sqli) ──
# case_id=TC-2add12cf
# case_name=[OWASP-API7] GET /api/diff — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/diff?from=%27+OR+1%3D1--
```json
null
```

HTTP 400

