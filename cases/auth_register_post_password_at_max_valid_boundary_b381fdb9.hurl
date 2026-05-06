# ── POST /auth/register - password at max_valid boundary ──
# case_id=TC-b381fdb9
# case_name=POST /auth/register - password at max_valid boundary
# step_id=step-main
# step_type=test
# technique=boundary_value
# priority=P1

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "kasandravelazquez@willis.org",
  "password": "zBlKzdHplyIohqMEAqvZSLUwRAAjdZKfbpkfEhUcSKoTKSlgMvwBEjoRpxXhryTaTAoTzCYyWaXpUkIgpumlAMpSEYEqFYHvmPDdtFumNUpHtbSoyugqaeiVyRdgqNwJsZzlXPJtrDBniDFcfYhHvlLEZBOqZCOoAPKPXTaHVHlRPRLPdCiRYyBYiVNGQIfRCXVbfVAECwwZbjBrGaKIfctBAjeidCzjvfjsjckVQIlqUrEHxrxTFDKxXvgrcFS"
}
```

HTTP 200

[Asserts]
duration < 2000

