# ══════════════════════════════════════════════════
# Create and retrieve user
# case_id=TC-8a91cfff
# case_name=Create and retrieve user
# case_kind=chain
# priority=P0
# ══════════════════════════════════════════════════

# ── POST new user [test] ────────────────────
# step_id=step-1
# step_type=test
# title=POST new user

POST {{base_url}}/users
```json
{
  "email": "jane.smith@example.com",
  "name": "Jane Smith"
}
```

HTTP 201

# ── GET created user by ID [test] ───────────
# step_id=step-2
# step_type=test
# title=GET created user by ID

GET {{base_url}}/users/1

HTTP 200

