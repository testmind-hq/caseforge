# ── [OWASP-API7] POST /api/admin/teams — injection (path-traversal) ──
# case_id=TC-a1f1c968
# case_name=[OWASP-API7] POST /api/admin/teams — injection (path-traversal)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "../../../etc/passwd"
}
```

HTTP 400

