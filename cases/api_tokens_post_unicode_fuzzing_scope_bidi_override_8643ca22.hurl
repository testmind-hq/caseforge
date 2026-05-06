# ── POST /api/tokens - [unicode_fuzzing] scope bidi_override ──
# case_id=TC-8643ca22
# case_name=POST /api/tokens - [unicode_fuzzing] scope bidi_override
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Amelia Cummings",
  "scope": "‮hello"
}
```

HTTP 400

