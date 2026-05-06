# ── [OWASP-API7] DELETE /api/admin/users/{id} — injection (sqli) ──
# case_id=TC-35704eb4
# case_name=[OWASP-API7] DELETE /api/admin/users/{id} — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

DELETE {{base_url}}/api/admin/users/%27%20OR%201=1--
```json
null
```

HTTP 400

