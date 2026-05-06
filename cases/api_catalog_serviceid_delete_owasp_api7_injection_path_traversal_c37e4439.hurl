# ── [OWASP-API7] DELETE /api/catalog/:serviceId — injection (path-traversal) ──
# case_id=TC-c37e4439
# case_name=[OWASP-API7] DELETE /api/catalog/:serviceId — injection (path-traversal)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

DELETE {{base_url}}/api/catalog/:serviceId
```json
null
```

HTTP 400

