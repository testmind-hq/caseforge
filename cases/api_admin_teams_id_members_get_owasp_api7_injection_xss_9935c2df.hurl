# ── [OWASP-API7] GET /api/admin/teams/{id}/members — injection (xss) ──
# case_id=TC-9935c2df
# case_name=[OWASP-API7] GET /api/admin/teams/{id}/members — injection (xss)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/admin/teams/%22%3E%3Cscript%3Ealert%281%29%3C%2Fscript%3E/members
```json
null
```

HTTP 400

