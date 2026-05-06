# ══════════════════════════════════════════════════
# Fail to create duplicate user with existing email
# case_id=TC-78c9e99f
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
  "email": "jane.doe@example.com",
  "name": "Jane Doe",
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
  "email": "jane.doe@example.com",
  "name": "Jane Smith",
  "password": "DifferentPass456!"
}
```

HTTP 409

