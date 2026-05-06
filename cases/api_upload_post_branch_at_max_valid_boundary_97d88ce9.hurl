# ── POST /api/upload - branch at max_valid boundary ──
# case_id=TC-97d88ce9
# case_name=POST /api/upload - branch at max_valid boundary
# step_id=step-main
# step_type=test
# technique=boundary_value
# priority=P1

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "shlwKqFxRFaVTdGNnBXhsNxUFKQKzOqqCpWDSXqaghfbdFJIOYfkDfFCtbwSekckstHPRyDaMVWZVWRBkbIgtUJDXhFeMmsQbiKempTLkISShAcAmWyGwOABgtbYCVEFRMDgKJWLKPmhAtLhMCfQaicCaLcxzIlibqzCyRCDxwtHNNlvPLxMHtmKcmYUtqMBHkdEiCZvhHNvCBGgJjhsNpbEGSpHxdHKXjeFulMWOPsstdqgeeJDWdLgyWSEFNF",
  "commitSha": "horde",
  "service": "patrol",
  "specContent": "early"
}
```

HTTP 200

[Asserts]
duration < 2000

