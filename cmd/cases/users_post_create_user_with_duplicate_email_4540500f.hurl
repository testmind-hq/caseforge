# ══════════════════════════════════════════════════
# Create user with duplicate email
# case_id=TC-4540500f
# case_name=Create user with duplicate email
# case_kind=chain
# priority=P1
# ══════════════════════════════════════════════════

# ── Create first user [test] ────────────────
# step_id=step-1
# step_type=test
# title=Create first user

POST {{base_url}}/users
```json
{
  "email": "alice@example.com",
  "name": "Alice Johnson",
  "password": "SecurePass123!"
}
```

HTTP 201

# ── Attempt to create second user with same email [test] ──
# step_id=step-2
# step_type=test
# title=Attempt to create second user with same email

POST {{base_url}}/users
```json
{
  "email": "alice@example.com",
  "name": "Alice Smith",
  "password": "DifferentPass456!"
}
```

HTTP 409

