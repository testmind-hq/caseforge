# ── [OWASP-API7] GET /api/specs/:service/versions — injection (sqli) ──
# case_id=TC-ffc707f5
# case_name=[OWASP-API7] GET /api/specs/:service/versions — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/specs/:service/versions
```json
null
```

HTTP 400

