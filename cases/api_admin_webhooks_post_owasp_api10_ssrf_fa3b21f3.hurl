# ── [OWASP-API10] POST /api/admin/webhooks — SSRF ──
# case_id=TC-fa3b21f3
# case_name=[OWASP-API10] POST /api/admin/webhooks — SSRF
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "url": "http://127.0.0.1"
}
```

HTTP 400

