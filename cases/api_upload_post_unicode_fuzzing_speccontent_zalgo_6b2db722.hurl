# ── POST /api/upload - [unicode_fuzzing] specContent zalgo ──
# case_id=TC-6b2db722
# case_name=POST /api/upload - [unicode_fuzzing] specContent zalgo
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
  "service": "consequently",
  "specContent": "z̀́̂̃̄̅̆̇a"
}
```

HTTP 400

