# ── [OWASP-API7] POST /api/admin/teams/{id}/grants — injection (path-traversal) ──
# case_id=TC-aa0b7128
# case_name=[OWASP-API7] POST /api/admin/teams/{id}/grants — injection (path-traversal)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/teams/..%2F..%2F..%2Fetc%2Fpasswd/grants
```json
null
```

HTTP 400

