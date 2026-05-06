# ── [OWASP-API7] DELETE /api/tokens/{id} — injection (path-traversal) ──
# case_id=TC-85b86fe3
# case_name=[OWASP-API7] DELETE /api/tokens/{id} — injection (path-traversal)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

DELETE {{base_url}}/api/tokens/..%2F..%2F..%2Fetc%2Fpasswd
```json
null
```

HTTP 400

