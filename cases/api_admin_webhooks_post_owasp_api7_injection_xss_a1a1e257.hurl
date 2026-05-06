# ── [OWASP-API7] POST /api/admin/webhooks — injection (xss) ──
# case_id=TC-a1a1e257
# case_name=[OWASP-API7] POST /api/admin/webhooks — injection (xss)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "providerType": "\"\u003e\u003cscript\u003ealert(1)\u003c/script\u003e"
}
```

HTTP 400

