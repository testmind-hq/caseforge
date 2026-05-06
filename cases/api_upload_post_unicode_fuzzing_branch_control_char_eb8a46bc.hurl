# ── POST /api/upload - [unicode_fuzzing] branch control_char ──
# case_id=TC-eb8a46bc
# case_name=POST /api/upload - [unicode_fuzzing] branch control_char
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "hello\u0000world",
  "commitSha": "herself",
  "service": "consequently",
  "specContent": "neither"
}
```

HTTP 400

