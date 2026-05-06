# ── [OWASP-API7] POST /auth/register — injection (xss) ──
# case_id=TC-368fd7b5
# case_name=[OWASP-API7] POST /auth/register — injection (xss)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "\"\u003e\u003cscript\u003ealert(1)\u003c/script\u003e"
}
```

HTTP 400

