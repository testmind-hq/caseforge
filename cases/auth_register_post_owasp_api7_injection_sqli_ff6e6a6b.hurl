# ── [OWASP-API7] POST /auth/register — injection (sqli) ──
# case_id=TC-ff6e6a6b
# case_name=[OWASP-API7] POST /auth/register — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "' OR 1=1--"
}
```

HTTP 400

