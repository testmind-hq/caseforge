# ── POST /api/upload - [unicode_fuzzing] service bidi_override ──
# case_id=TC-71d03103
# case_name=POST /api/upload - [unicode_fuzzing] service bidi_override
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
  "service": "‮hello",
  "specContent": "neither"
}
```

HTTP 400

