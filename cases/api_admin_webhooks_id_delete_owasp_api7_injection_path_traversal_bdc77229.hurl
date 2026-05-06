# ── [OWASP-API7] DELETE /api/admin/webhooks/:id — injection (path-traversal) ──
# case_id=TC-bdc77229
# case_name=[OWASP-API7] DELETE /api/admin/webhooks/:id — injection (path-traversal)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

DELETE {{base_url}}/api/admin/webhooks/:id
```json
null
```

HTTP 400

