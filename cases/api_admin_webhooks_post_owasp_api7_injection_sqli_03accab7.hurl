# ── [OWASP-API7] POST /api/admin/webhooks — injection (sqli) ──
# case_id=TC-03accab7
# case_name=[OWASP-API7] POST /api/admin/webhooks — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "providerType": "' OR 1=1--"
}
```

HTTP 400

