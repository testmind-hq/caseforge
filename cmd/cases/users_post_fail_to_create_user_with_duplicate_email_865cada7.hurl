# ══════════════════════════════════════════════════
# Fail to create user with duplicate email
# case_id=TC-865cada7
# case_name=Fail to create user with duplicate email
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
  "email": "john.doe@example.com",
  "name": "John Doe",
  "password": "SecurePass123!"
}
```

HTTP 201

# ── Attempt to create user with same email [test] ──
# step_id=step-2
# step_type=test
# title=Attempt to create user with same email

POST {{base_url}}/users
```json
{
  "email": "john.doe@example.com",
  "name": "Another John",
  "password": "DifferentPass456!"
}
```

HTTP 409

