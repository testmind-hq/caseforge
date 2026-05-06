# ── [OWASP-API7] POST /api/admin/teams — injection (sqli) ──
# case_id=TC-3e99ea9b
# case_name=[OWASP-API7] POST /api/admin/teams — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "' OR 1=1--"
}
```

HTTP 400

