# ── POST /api/upload - [unicode_fuzzing] branch zalgo ──
# case_id=TC-3c16d4b3
# case_name=POST /api/upload - [unicode_fuzzing] branch zalgo
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "z̀́̂̃̄̅̆̇a",
  "commitSha": "herself",
  "service": "consequently",
  "specContent": "neither"
}
```

HTTP 400

