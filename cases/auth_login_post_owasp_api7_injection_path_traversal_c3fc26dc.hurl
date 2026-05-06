# ── [OWASP-API7] POST /auth/login — injection (path-traversal) ──
# case_id=TC-c3fc26dc
# case_name=[OWASP-API7] POST /auth/login — injection (path-traversal)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "../../../etc/passwd"
}
```

HTTP 400

