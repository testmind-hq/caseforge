# ── [OWASP-API7] POST /api/admin/teams/{id}/grants — injection (sqli) ──
# case_id=TC-ea6fd919
# case_name=[OWASP-API7] POST /api/admin/teams/{id}/grants — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/teams/%27%20OR%201=1--/grants
```json
null
```

HTTP 400

