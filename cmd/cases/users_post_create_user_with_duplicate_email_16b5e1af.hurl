# ── Create user with duplicate email ────────
# case_id=TC-16b5e1af
# case_name=Create user with duplicate email
# step_id=step-1
# step_type=test
# technique=ask
# priority=P1

POST {{base_url}}/users
```json
{
  "email": "existing.user@example.com",
  "name": "Jane Doe",
  "password": "SecurePass123!"
}
```

HTTP 409

