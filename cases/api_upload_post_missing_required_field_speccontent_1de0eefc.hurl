# ── POST /api/upload - missing required field "specContent" ──
# case_id=TC-1de0eefc
# case_name=POST /api/upload - missing required field "specContent"
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P1

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "can",
  "commitSha": "why",
  "service": "forest"
}
```

HTTP 422

