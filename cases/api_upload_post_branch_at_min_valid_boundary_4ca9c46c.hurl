# ── POST /api/upload - branch at min_valid boundary ──
# case_id=TC-4ca9c46c
# case_name=POST /api/upload - branch at min_valid boundary
# step_id=step-main
# step_type=test
# technique=boundary_value
# priority=P1

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "b",
  "commitSha": "horde",
  "service": "patrol",
  "specContent": "early"
}
```

HTTP 200

[Asserts]
duration < 2000

