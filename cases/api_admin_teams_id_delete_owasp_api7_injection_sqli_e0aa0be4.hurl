# ── [OWASP-API7] DELETE /api/admin/teams/{id} — injection (sqli) ──
# case_id=TC-e0aa0be4
# case_name=[OWASP-API7] DELETE /api/admin/teams/{id} — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

DELETE {{base_url}}/api/admin/teams/%27%20OR%201=1--
```json
null
```

HTTP 400

