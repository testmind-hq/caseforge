# ── POST /auth/register - password at max_plus_one_invalid boundary ──
# case_id=TC-0de23fb9
# case_name=POST /auth/register - password at max_plus_one_invalid boundary
# step_id=step-main
# step_type=test
# technique=boundary_value
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "kasandravelazquez@willis.org",
  "password": "rPNlcdUMPwImsPdHFstXXMFIWbajRRdQloozwcKtoDbGhjiVVjHhIxcPpxMVGqqKfZycxZGoowdemLuYWOaEvFeerqBahGZywYIkuGXZrJdCNLryEunbqPYCHWypnUwNviWToCVJFisKyZtCteizZYgpdPlJDBzSucWfdtYFBAzmlDrKirFlAXDxVwWdZscUXFIAryQbydibyCuTJuKPjVPFBgydzlVHJwlOmkfnmyWhxdOnhlOMZdXVRggOpqya"
}
```

HTTP 422

