# ── [OWASP-API7] POST /api/upload — injection (xss) ──
# case_id=TC-81a2a747
# case_name=[OWASP-API7] POST /api/upload — injection (xss)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "\"\u003e\u003cscript\u003ealert(1)\u003c/script\u003e"
}
```

HTTP 400

