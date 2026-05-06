# ── [OWASP-API7] POST /api/admin/teams — injection (xss) ──
# case_id=TC-a582e336
# case_name=[OWASP-API7] POST /api/admin/teams — injection (xss)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "\"\u003e\u003cscript\u003ealert(1)\u003c/script\u003e"
}
```

HTTP 400

