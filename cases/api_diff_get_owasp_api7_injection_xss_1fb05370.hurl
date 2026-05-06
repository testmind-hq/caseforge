# ── [OWASP-API7] GET /api/diff — injection (xss) ──
# case_id=TC-1fb05370
# case_name=[OWASP-API7] GET /api/diff — injection (xss)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/diff?from=%22%3E%3Cscript%3Ealert%281%29%3C%2Fscript%3E
```json
null
```

HTTP 400

