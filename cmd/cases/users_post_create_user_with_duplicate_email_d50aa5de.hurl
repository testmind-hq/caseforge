# ══════════════════════════════════════════════════
# Create user with duplicate email
# case_id=TC-d50aa5de
# case_name=Create user with duplicate email
# case_kind=chain
# priority=P1
# ══════════════════════════════════════════════════

# ── POST first user [test] ──────────────────
# step_id=step-1
# step_type=test
# title=POST first user

POST {{base_url}}/users
```json
{
  "email": "jane.doe@example.com",
  "name": "Jane Doe",
  "password": "SecurePass123!"
}
```

HTTP 201

# ── POST second user with same email [test] ──
# step_id=step-2
# step_type=test
# title=POST second user with same email

POST {{base_url}}/users
```json
{
  "email": "jane.doe@example.com",
  "name": "Jane Smith",
  "password": "AnotherPass456!"
}
```

HTTP 409

