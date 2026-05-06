# ── [OWASP-API7] POST /api/admin/webhooks/:id/test — injection (xss) ──
# case_id=TC-e8743ba7
# case_name=[OWASP-API7] POST /api/admin/webhooks/:id/test — injection (xss)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/webhooks/:id/test
```json
null
```

HTTP 400

