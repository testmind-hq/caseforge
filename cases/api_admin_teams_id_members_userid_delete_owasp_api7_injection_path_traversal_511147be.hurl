# ── [OWASP-API7] DELETE /api/admin/teams/{id}/members/{userId} — injection (path-traversal) ──
# case_id=TC-511147be
# case_name=[OWASP-API7] DELETE /api/admin/teams/{id}/members/{userId} — injection (path-traversal)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

DELETE {{base_url}}/api/admin/teams/..%2F..%2F..%2Fetc%2Fpasswd/members/{userId}
```json
null
```

HTTP 400

