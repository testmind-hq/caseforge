# ── POST /api/upload - [unicode_fuzzing] branch bidi_override ──
# case_id=TC-09b46ba6
# case_name=POST /api/upload - [unicode_fuzzing] branch bidi_override
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "‮hello",
  "commitSha": "herself",
  "service": "consequently",
  "specContent": "neither"
}
```

HTTP 400

