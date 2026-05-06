# ── POST /api/tokens - [unicode_fuzzing] name bidi_override ──
# case_id=TC-33a5a9d7
# case_name=POST /api/tokens - [unicode_fuzzing] name bidi_override
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "‮hello",
  "scope": "read"
}
```

HTTP 400

