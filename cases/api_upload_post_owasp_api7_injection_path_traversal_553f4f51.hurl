# ── [OWASP-API7] POST /api/upload — injection (path-traversal) ──
# case_id=TC-553f4f51
# case_name=[OWASP-API7] POST /api/upload — injection (path-traversal)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "../../../etc/passwd"
}
```

HTTP 400

