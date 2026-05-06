# ── [OWASP-API7] GET /api/specs/:service/versions — injection (path-traversal) ──
# case_id=TC-106c80c0
# case_name=[OWASP-API7] GET /api/specs/:service/versions — injection (path-traversal)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/specs/:service/versions
```json
null
```

HTTP 400

