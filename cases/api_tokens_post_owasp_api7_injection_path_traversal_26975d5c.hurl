# ── [OWASP-API7] POST /api/tokens — injection (path-traversal) ──
# case_id=TC-26975d5c
# case_name=[OWASP-API7] POST /api/tokens — injection (path-traversal)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "scope": "../../../etc/passwd"
}
```

HTTP 400

