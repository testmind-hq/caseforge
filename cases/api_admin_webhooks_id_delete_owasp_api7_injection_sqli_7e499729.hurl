# ── [OWASP-API7] DELETE /api/admin/webhooks/:id — injection (sqli) ──
# case_id=TC-7e499729
# case_name=[OWASP-API7] DELETE /api/admin/webhooks/:id — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

DELETE {{base_url}}/api/admin/webhooks/:id
```json
null
```

HTTP 400

