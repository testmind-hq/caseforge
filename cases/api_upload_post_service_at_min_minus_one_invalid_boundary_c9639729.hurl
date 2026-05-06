# ── POST /api/upload - service at min_minus_one_invalid boundary ──
# case_id=TC-c9639729
# case_name=POST /api/upload - service at min_minus_one_invalid boundary
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
  "service": "P",
  "specContent": "early"
}
```

HTTP 422

