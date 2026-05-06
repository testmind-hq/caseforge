# ── [OWASP-API7] PUT /api/admin/users/{id} — injection (sqli) ──
# case_id=TC-c653b26d
# case_name=[OWASP-API7] PUT /api/admin/users/{id} — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PUT {{base_url}}/api/admin/users/%27%20OR%201=1--
```json
null
```

HTTP 400

