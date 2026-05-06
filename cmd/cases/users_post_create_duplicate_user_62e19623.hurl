# ══════════════════════════════════════════════════
# Create duplicate user
# case_id=TC-62e19623
# case_name=Create duplicate user
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

# ── Attempt to create duplicate user [test] ──
# step_id=step-2
# step_type=test
# title=Attempt to create duplicate user

POST {{base_url}}/users
```json
{
  "email": "jane.doe@example.com",
  "name": "Jane Doe",
  "password": "AnotherPass456!"
}
```

HTTP 409

