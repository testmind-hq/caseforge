# ── POST /api/upload - [unicode_fuzzing] service zalgo ──
# case_id=TC-7d8cc30e
# case_name=POST /api/upload - [unicode_fuzzing] service zalgo
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "honestly",
  "commitSha": "herself",
  "service": "z̀́̂̃̄̅̆̇a",
  "specContent": "neither"
}
```

HTTP 400

