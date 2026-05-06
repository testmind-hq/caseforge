# ── [OWASP-API7] DELETE /api/admin/teams/{id} — injection (xss) ──
# case_id=TC-cdcba009
# case_name=[OWASP-API7] DELETE /api/admin/teams/{id} — injection (xss)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

DELETE {{base_url}}/api/admin/teams/%22%3E%3Cscript%3Ealert%281%29%3C%2Fscript%3E
```json
null
```

HTTP 400

