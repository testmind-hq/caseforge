# ── POST /api/upload - [unicode_fuzzing] branch zero_width ──
# case_id=TC-d4d96d5e
# case_name=POST /api/upload - [unicode_fuzzing] branch zero_width
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "​hello",
  "commitSha": "herself",
  "service": "consequently",
  "specContent": "neither"
}
```

HTTP 400

