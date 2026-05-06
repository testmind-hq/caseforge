# ══════════════════════════════════════════════════
# Fail to create duplicate user with existing email
# case_id=TC-b9e88eb8
# case_name=Fail to create duplicate user with existing email
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

# ── Attempt to create user with same email [test] ──
# step_id=step-2
# step_type=test
# title=Attempt to create user with same email

POST {{base_url}}/users
```json
{
  "email": "alice@example.com",
  "name": "Alice Johnson",
  "password": "DifferentPass456!"
}
```

HTTP 409

