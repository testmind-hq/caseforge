# ── POST /api/upload - missing required field "service" ──
# case_id=TC-89850cfa
# case_name=POST /api/upload - missing required field "service"
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P1

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "occasionally",
  "commitSha": "lastly",
  "specContent": "eat"
}
```

HTTP 422

