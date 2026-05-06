# ── POST /auth/login - wrong content-type (text/plain) ──
# case_id=TC-ea0be7b9
# case_name=POST /auth/login - wrong content-type (text/plain)
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/auth/login
Content-Type: text/plain
```json
{
  "email": "ottonorris@sullivan.com",
  "password": "float"
}
```

HTTP 415

