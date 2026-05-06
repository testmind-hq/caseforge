# ── POST /api/tokens - [unicode_fuzzing] name zalgo ──
# case_id=TC-431d2bbf
# case_name=POST /api/tokens - [unicode_fuzzing] name zalgo
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "z̀́̂̃̄̅̆̇a",
  "scope": "read"
}
```

HTTP 400

