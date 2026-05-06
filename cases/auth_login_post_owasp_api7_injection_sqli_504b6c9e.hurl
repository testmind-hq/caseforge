# ── [OWASP-API7] POST /auth/login — injection (sqli) ──
# case_id=TC-504b6c9e
# case_name=[OWASP-API7] POST /auth/login — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "' OR 1=1--"
}
```

HTTP 400

