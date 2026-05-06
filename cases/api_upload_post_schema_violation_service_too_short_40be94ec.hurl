# ── POST /api/upload - [schema_violation] service_too_short ──
# case_id=TC-40be94ec
# case_name=POST /api/upload - [schema_violation] service_too_short
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "to",
  "commitSha": "Brazilian",
  "service": "",
  "specContent": "tonight"
}
```

HTTP 422

