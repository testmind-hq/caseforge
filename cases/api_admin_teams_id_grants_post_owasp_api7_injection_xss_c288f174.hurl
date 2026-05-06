# ── [OWASP-API7] POST /api/admin/teams/{id}/grants — injection (xss) ──
# case_id=TC-c288f174
# case_name=[OWASP-API7] POST /api/admin/teams/{id}/grants — injection (xss)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/teams/%22%3E%3Cscript%3Ealert%281%29%3C%2Fscript%3E/grants
```json
null
```

HTTP 400

