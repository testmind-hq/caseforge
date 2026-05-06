# ── [OWASP-API7] GET /api/admin/teams/{id}/services — injection (sqli) ──
# case_id=TC-2e72efb4
# case_name=[OWASP-API7] GET /api/admin/teams/{id}/services — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/admin/teams/%27%20OR%201=1--/services
```json
null
```

HTTP 400

