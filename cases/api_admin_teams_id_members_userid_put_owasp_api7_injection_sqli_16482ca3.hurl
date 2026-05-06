# ── [OWASP-API7] PUT /api/admin/teams/{id}/members/{userId} — injection (sqli) ──
# case_id=TC-16482ca3
# case_name=[OWASP-API7] PUT /api/admin/teams/{id}/members/{userId} — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PUT {{base_url}}/api/admin/teams/%27%20OR%201=1--/members/{userId}
```json
null
```

HTTP 400

