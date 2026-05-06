# ── [OWASP-API7] POST /api/upload — injection (sqli) ──
# case_id=TC-b528a6e6
# case_name=[OWASP-API7] POST /api/upload — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "' OR 1=1--"
}
```

HTTP 400

