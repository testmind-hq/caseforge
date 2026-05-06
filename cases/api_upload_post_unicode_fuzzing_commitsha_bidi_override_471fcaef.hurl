# ── POST /api/upload - [unicode_fuzzing] commitSha bidi_override ──
# case_id=TC-471fcaef
# case_name=POST /api/upload - [unicode_fuzzing] commitSha bidi_override
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "honestly",
  "commitSha": "‮hello",
  "service": "consequently",
  "specContent": "neither"
}
```

HTTP 400

