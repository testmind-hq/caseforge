# ── POST /auth/register - wrong content-type (text/plain) ──
# case_id=TC-9cf203de
# case_name=POST /auth/register - wrong content-type (text/plain)
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/auth/register
Content-Type: text/plain
```json
{
  "email": "audreygarrett@morris.info",
  "password": "mouth"
}
```

HTTP 415

