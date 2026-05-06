# ── POST /api/upload - [unicode_fuzzing] commitSha control_char ──
# case_id=TC-1e3b28af
# case_name=POST /api/upload - [unicode_fuzzing] commitSha control_char
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "honestly",
  "commitSha": "hello\u0000world",
  "service": "consequently",
  "specContent": "neither"
}
```

HTTP 400

