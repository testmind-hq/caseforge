# ── [OWASP-API7] GET /api/specs/:service/versions — injection (xss) ──
# case_id=TC-cf42e9f4
# case_name=[OWASP-API7] GET /api/specs/:service/versions — injection (xss)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/specs/:service/versions
```json
null
```

HTTP 400

