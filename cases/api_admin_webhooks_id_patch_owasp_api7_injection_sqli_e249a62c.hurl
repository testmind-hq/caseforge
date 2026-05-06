# ── [OWASP-API7] PATCH /api/admin/webhooks/:id — injection (sqli) ──
# case_id=TC-e249a62c
# case_name=[OWASP-API7] PATCH /api/admin/webhooks/:id — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PATCH {{base_url}}/api/admin/webhooks/:id
```json
null
```

HTTP 400

