# ── POST /api/tokens - name at max_plus_one_invalid boundary ──
# case_id=TC-7b3217ba
# case_name=POST /api/tokens - name at max_plus_one_invalid boundary
# step_id=step-main
# step_type=test
# technique=boundary_value
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "gZkkKaBcgUUrJhMvMmXsjgUJDOfrVpkfGCKVAUujjHuMbmjqYrroOdpRDCHXNKftgwkIjzdVDnyjNbwYqqZrajsqPvSTaCwhMFwMjAZyBQIjmghcfkelirBpAPxhbuYkwsodExCcRneWXSlyLvtcufLRHJWucpZNlpPiKuSLlicpZPdObnVxJdhXykuHmqCapfBevaSSFSPEtYlzUlPAVbisIBFXneKSEoFFcgPCMSeUhOCBMxaqhfiLFJvQwWsX",
  "scope": "read"
}
```

HTTP 422

