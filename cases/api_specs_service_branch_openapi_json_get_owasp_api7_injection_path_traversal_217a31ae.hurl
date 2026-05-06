# ── [OWASP-API7] GET /api/specs/{service}/{branch}/openapi.json — injection (path-traversal) ──
# case_id=TC-217a31ae
# case_name=[OWASP-API7] GET /api/specs/{service}/{branch}/openapi.json — injection (path-traversal)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/specs/..%2F..%2F..%2Fetc%2Fpasswd/{branch}/openapi.json
```json
null
```

HTTP 400

