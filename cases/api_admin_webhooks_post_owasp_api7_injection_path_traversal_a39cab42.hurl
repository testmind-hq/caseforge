# ── [OWASP-API7] POST /api/admin/webhooks — injection (path-traversal) ──
# case_id=TC-a39cab42
# case_name=[OWASP-API7] POST /api/admin/webhooks — injection (path-traversal)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "providerType": "../../../etc/passwd"
}
```

HTTP 400

