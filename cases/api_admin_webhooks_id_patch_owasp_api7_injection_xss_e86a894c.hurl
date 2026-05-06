# ── [OWASP-API7] PATCH /api/admin/webhooks/:id — injection (xss) ──
# case_id=TC-e86a894c
# case_name=[OWASP-API7] PATCH /api/admin/webhooks/:id — injection (xss)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PATCH {{base_url}}/api/admin/webhooks/:id
```json
null
```

HTTP 400

