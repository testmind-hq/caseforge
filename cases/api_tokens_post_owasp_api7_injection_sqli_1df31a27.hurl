# ── [OWASP-API7] POST /api/tokens — injection (sqli) ──
# case_id=TC-1df31a27
# case_name=[OWASP-API7] POST /api/tokens — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "scope": "' OR 1=1--"
}
```

HTTP 400

