# ── POST /api/upload - [unicode_fuzzing] specContent control_char ──
# case_id=TC-7ff8ca85
# case_name=POST /api/upload - [unicode_fuzzing] specContent control_char
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
  "specContent": "hello\u0000world"
}
```

HTTP 400

