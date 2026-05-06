# ══════════════════════════════════════════════════
# Create user and retrieve it
# case_id=TC-f9ba7a73
# case_name=Create user and retrieve it
# case_kind=chain
# priority=P0
# ══════════════════════════════════════════════════

# ── Create new user [test] ──────────────────
# step_id=step-1
# step_type=test
# title=Create new user

POST {{base_url}}/users
```json
{
  "email": "alice.smith@example.com",
  "name": "Alice Smith",
  "password": "SecurePass123!"
}
```

HTTP 201

# ── Retrieve created user [test] ────────────
# step_id=step-2
# step_type=test
# title=Retrieve created user

GET {{base_url}}/users/1

HTTP 200

