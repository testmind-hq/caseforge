# ── POST /api/admin/teams - [unicode_fuzzing] name bidi_override ──
# case_id=TC-19447855
# case_name=POST /api/admin/teams - [unicode_fuzzing] name bidi_override
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Review the woman every 2 weeks.",
  "displayName": "example",
  "name": "‮hello"
}
```

HTTP 400

