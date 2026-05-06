# ── [OWASP-API7] PUT /api/admin/teams/{id} — injection (sqli) ──
# case_id=TC-c7f786e4
# case_name=[OWASP-API7] PUT /api/admin/teams/{id} — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PUT {{base_url}}/api/admin/teams/%27%20OR%201=1--
```json
null
```

HTTP 400

