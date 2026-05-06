# ══════════════════════════════════════════════════
# Create user with duplicate email
# case_id=TC-2143a276
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
  "email": "duplicate@example.com",
  "name": "First User",
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
  "email": "duplicate@example.com",
  "name": "Second User",
  "password": "DifferentPass456!"
}
```

HTTP 409

