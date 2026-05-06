# ── POST /api/upload - [unicode_fuzzing] specContent bidi_override ──
# case_id=TC-131ad5f4
# case_name=POST /api/upload - [unicode_fuzzing] specContent bidi_override
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
  "specContent": "‮hello"
}
```

HTTP 400

