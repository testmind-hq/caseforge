# ── POST /api/upload - [schema_violation] specContent_too_short ──
# case_id=TC-af512611
# case_name=POST /api/upload - [schema_violation] specContent_too_short
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
  "service": "intimidate",
  "specContent": ""
}
```

HTTP 422

