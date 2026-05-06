# ── [OWASP-API7] PUT /api/admin/services/{serviceId}/team — injection (sqli) ──
# case_id=TC-53f0e55f
# case_name=[OWASP-API7] PUT /api/admin/services/{serviceId}/team — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PUT {{base_url}}/api/admin/services/%27%20OR%201=1--/team
```json
null
```

HTTP 400

