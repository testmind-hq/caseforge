# ── [OWASP-API7] GET /api/specs/{service}/{branch}/openapi.json — injection (sqli) ──
# case_id=TC-3e62652b
# case_name=[OWASP-API7] GET /api/specs/{service}/{branch}/openapi.json — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/specs/%27%20OR%201=1--/{branch}/openapi.json
```json
null
```

HTTP 400

