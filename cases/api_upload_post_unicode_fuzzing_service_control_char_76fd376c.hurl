# ── POST /api/upload - [unicode_fuzzing] service control_char ──
# case_id=TC-76fd376c
# case_name=POST /api/upload - [unicode_fuzzing] service control_char
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
  "service": "hello\u0000world",
  "specContent": "neither"
}
```

HTTP 400

