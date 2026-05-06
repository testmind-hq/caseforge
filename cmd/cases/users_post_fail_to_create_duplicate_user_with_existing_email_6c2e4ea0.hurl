# ══════════════════════════════════════════════════
# Fail to create duplicate user with existing email
# case_id=TC-6c2e4ea0
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
  "email": "alice.brown@example.com",
  "name": "Alice Brown",
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
  "email": "alice.brown@example.com",
  "name": "Alice Brown",
  "password": "DifferentPass456!"
}
```

HTTP 409

