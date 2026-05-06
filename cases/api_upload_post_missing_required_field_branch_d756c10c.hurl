# ── POST /api/upload - missing required field "branch" ──
# case_id=TC-d756c10c
# case_name=POST /api/upload - missing required field "branch"
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P1

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "commitSha": "news",
  "service": "seldom",
  "specContent": "who"
}
```

HTTP 422

