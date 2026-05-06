# ── [OWASP-API7] PUT /api/admin/services/{serviceId}/team — injection (xss) ──
# case_id=TC-3ad867af
# case_name=[OWASP-API7] PUT /api/admin/services/{serviceId}/team — injection (xss)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PUT {{base_url}}/api/admin/services/%22%3E%3Cscript%3Ealert%281%29%3C%2Fscript%3E/team
```json
null
```

HTTP 400

