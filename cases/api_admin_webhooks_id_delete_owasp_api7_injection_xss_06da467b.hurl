# ── [OWASP-API7] DELETE /api/admin/webhooks/:id — injection (xss) ──
# case_id=TC-06da467b
# case_name=[OWASP-API7] DELETE /api/admin/webhooks/:id — injection (xss)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

DELETE {{base_url}}/api/admin/webhooks/:id
```json
null
```

HTTP 400

