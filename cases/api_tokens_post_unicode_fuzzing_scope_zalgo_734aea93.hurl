# ── POST /api/tokens - [unicode_fuzzing] scope zalgo ──
# case_id=TC-734aea93
# case_name=POST /api/tokens - [unicode_fuzzing] scope zalgo
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Amelia Cummings",
  "scope": "z̀́̂̃̄̅̆̇a"
}
```

HTTP 400

