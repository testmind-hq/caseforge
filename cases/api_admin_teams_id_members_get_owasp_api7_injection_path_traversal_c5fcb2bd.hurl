# ── [OWASP-API7] GET /api/admin/teams/{id}/members — injection (path-traversal) ──
# case_id=TC-c5fcb2bd
# case_name=[OWASP-API7] GET /api/admin/teams/{id}/members — injection (path-traversal)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

GET {{base_url}}/api/admin/teams/..%2F..%2F..%2Fetc%2Fpasswd/members
```json
null
```

HTTP 400

