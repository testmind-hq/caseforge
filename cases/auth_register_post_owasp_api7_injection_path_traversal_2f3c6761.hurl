# ── [OWASP-API7] POST /auth/register — injection (path-traversal) ──
# case_id=TC-2f3c6761
# case_name=[OWASP-API7] POST /auth/register — injection (path-traversal)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "../../../etc/passwd"
}
```

HTTP 400

