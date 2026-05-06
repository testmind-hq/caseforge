# ── [OWASP-API7] GET /api/diff — injection (path-traversal) ──
# case_id=TC-d2e88748
# case_name=[OWASP-API7] GET /api/diff — injection (path-traversal)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/diff?from=..%2F..%2F..%2Fetc%2Fpasswd
```json
null
```

HTTP 400

