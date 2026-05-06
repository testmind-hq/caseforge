# ── POST /api/upload - service at min_valid boundary ──
# case_id=TC-fa5f2879
# case_name=POST /api/upload - service at min_valid boundary
# step_id=step-main
# step_type=test
# technique=boundary_value
# priority=P1

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "annoying",
  "commitSha": "horde",
  "service": "v",
  "specContent": "early"
}
```

HTTP 200

[Asserts]
duration < 2000

