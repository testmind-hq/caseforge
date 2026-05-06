# ── POST /api/upload - [unicode_fuzzing] service zero_width ──
# case_id=TC-f8f99bf7
# case_name=POST /api/upload - [unicode_fuzzing] service zero_width
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
  "service": "​hello",
  "specContent": "neither"
}
```

HTTP 400

