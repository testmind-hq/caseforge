# ── POST /api/upload - [unicode_fuzzing] specContent zero_width ──
# case_id=TC-7ac120c3
# case_name=POST /api/upload - [unicode_fuzzing] specContent zero_width
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
  "specContent": "​hello"
}
```

HTTP 400

