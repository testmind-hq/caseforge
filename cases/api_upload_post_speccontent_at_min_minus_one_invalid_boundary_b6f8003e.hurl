# ── POST /api/upload - specContent at min_minus_one_invalid boundary ──
# case_id=TC-b6f8003e
# case_name=POST /api/upload - specContent at min_minus_one_invalid boundary
# step_id=step-main
# step_type=test
# technique=boundary_value
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "annoying",
  "commitSha": "horde",
  "service": "patrol",
  "specContent": "E"
}
```

HTTP 422

