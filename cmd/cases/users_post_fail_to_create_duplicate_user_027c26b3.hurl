# ══════════════════════════════════════════════════
# Fail to create duplicate user
# case_id=TC-027c26b3
# case_name=Fail to create duplicate user
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
  "name": "Duplicate Test"
}
```

HTTP 201

# ── POST duplicate user with same email [test] ──
# step_id=step-2
# step_type=test
# title=POST duplicate user with same email

POST {{base_url}}/users
```json
{
  "email": "duplicate@example.com",
  "name": "Duplicate Test"
}
```

HTTP 409

