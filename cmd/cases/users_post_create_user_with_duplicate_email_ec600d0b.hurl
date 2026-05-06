# ══════════════════════════════════════════════════
# Create user with duplicate email
# case_id=TC-ec600d0b
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
  "email": "alice.brown@example.com",
  "name": "Alice Brown",
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
  "email": "alice.brown@example.com",
  "name": "Alice Duplicate",
  "password": "AnotherPass456!"
}
```

HTTP 409

