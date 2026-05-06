# ══════════════════════════════════════════════════
# Create duplicate user with existing email
# case_id=TC-7c11147b
# case_name=Create duplicate user with existing email
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
  "email": "duplicate@example.com",
  "name": "Duplicate User",
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
  "email": "duplicate@example.com",
  "name": "Another User",
  "password": "DifferentPass456!"
}
```

HTTP 409

