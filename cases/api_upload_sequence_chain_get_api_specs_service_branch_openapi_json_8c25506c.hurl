# ══════════════════════════════════════════════════
# sequence chain: /api/upload → GET /api/specs/{service}/{branch}/openapi.json
# case_id=TC-8c25506c
# case_name=sequence chain: /api/upload → GET /api/specs/{service}/{branch}/openapi.json
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
  "branch": "gift",
  "commitSha": "host",
  "service": "been",
  "specContent": "time"
}
```

HTTP *

[Captures]
service: jsonpath "$.service"

[Asserts]
status < 300

# ── use via GET /api/specs/{service}/{branch}/openapi.json [test] ──
# step_id=step-test
# step_type=test
# title=use via GET /api/specs/{service}/{branch}/openapi.json
# depends_on=step-setup

GET {{base_url}}/api/specs/{{service}}/{branch}/openapi.json

HTTP *

[Asserts]
status < 300

