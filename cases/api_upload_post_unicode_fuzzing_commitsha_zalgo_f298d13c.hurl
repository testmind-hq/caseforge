# ── POST /api/upload - [unicode_fuzzing] commitSha zalgo ──
# case_id=TC-f298d13c
# case_name=POST /api/upload - [unicode_fuzzing] commitSha zalgo
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "honestly",
  "commitSha": "z̀́̂̃̄̅̆̇a",
  "service": "consequently",
  "specContent": "neither"
}
```

HTTP 400

