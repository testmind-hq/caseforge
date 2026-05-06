# ── POST /api/upload - specContent at min_valid boundary ──
# case_id=TC-edc8ded2
# case_name=POST /api/upload - specContent at min_valid boundary
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
  "service": "patrol",
  "specContent": "s"
}
```

HTTP 200

[Asserts]
duration < 2000

