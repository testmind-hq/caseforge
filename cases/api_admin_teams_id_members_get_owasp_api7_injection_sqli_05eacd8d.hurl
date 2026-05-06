# ── [OWASP-API7] GET /api/admin/teams/{id}/members — injection (sqli) ──
# case_id=TC-05eacd8d
# case_name=[OWASP-API7] GET /api/admin/teams/{id}/members — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/admin/teams/%27%20OR%201=1--/members
```json
null
```

HTTP 400

