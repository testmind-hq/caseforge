# ══════════════════════════════════════════════════
# Create user with duplicate email
# case_id=TC-0be9ec08
# case_name=Create user with duplicate email
# case_kind=chain
# priority=P1
# ══════════════════════════════════════════════════

# ── POST first user successfully [test] ─────
# step_id=step-1
# step_type=test
# title=POST first user successfully

POST {{base_url}}/users
```json
{
  "email": "alice@example.com",
  "name": "Alice Johnson",
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
  "email": "alice@example.com",
  "name": "Alice Smith",
  "password": "DifferentPass456!"
}
```

HTTP 409

