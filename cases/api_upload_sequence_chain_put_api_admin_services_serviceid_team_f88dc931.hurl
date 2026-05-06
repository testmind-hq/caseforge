# ══════════════════════════════════════════════════
# sequence chain: /api/upload → PUT /api/admin/services/{serviceId}/team
# case_id=TC-f88dc931
# case_name=sequence chain: /api/upload → PUT /api/admin/services/{serviceId}/team
# case_kind=chain
# priority=P1
# ══════════════════════════════════════════════════

# ── create via POST /api/upload [setup] ─────
# step_id=step-setup
# step_type=setup
# title=create via POST /api/upload

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "someone",
  "commitSha": "instead",
  "service": "therefore",
  "specContent": "yesterday"
}
```

HTTP *

[Captures]
serviceId: jsonpath "$.service"

[Asserts]
status < 300

# ── use via PUT /api/admin/services/{serviceId}/team [test] ──
# step_id=step-test
# step_type=test
# title=use via PUT /api/admin/services/{serviceId}/team
# depends_on=step-setup

PUT {{base_url}}/api/admin/services/{{serviceId}}/team
Content-Type: application/json
```json
{
  "teamId": "e76c96fd-19bb-41c3-a5a4-6720d313f439"
}
```

HTTP *

[Asserts]
status < 300

