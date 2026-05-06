# ── [OWASP-API7] PUT /api/admin/users/{id} — injection (xss) ──
# case_id=TC-51b9a625
# case_name=[OWASP-API7] PUT /api/admin/users/{id} — injection (xss)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PUT {{base_url}}/api/admin/users/%22%3E%3Cscript%3Ealert%281%29%3C%2Fscript%3E
```json
null
```

HTTP 400

