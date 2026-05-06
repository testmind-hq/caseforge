# ── [OWASP-API7] GET /api/specs/{service}/{branch}/openapi.json — injection (xss) ──
# case_id=TC-69cf35a6
# case_name=[OWASP-API7] GET /api/specs/{service}/{branch}/openapi.json — injection (xss)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/specs/%22%3E%3Cscript%3Ealert%281%29%3C%2Fscript%3E/{branch}/openapi.json
```json
null
```

HTTP 400

