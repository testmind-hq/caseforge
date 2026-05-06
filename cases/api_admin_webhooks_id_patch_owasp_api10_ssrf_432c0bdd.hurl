# ── [OWASP-API10] PATCH /api/admin/webhooks/:id — SSRF ──
# case_id=TC-432c0bdd
# case_name=[OWASP-API10] PATCH /api/admin/webhooks/:id — SSRF
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: application/json
```json
{
  "url": "http://127.0.0.1"
}
```

HTTP 400

