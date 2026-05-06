# ── [OWASP-API7] DELETE /api/catalog/:serviceId — injection (xss) ──
# case_id=TC-bfdae539
# case_name=[OWASP-API7] DELETE /api/catalog/:serviceId — injection (xss)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

DELETE {{base_url}}/api/catalog/:serviceId
```json
null
```

HTTP 400

