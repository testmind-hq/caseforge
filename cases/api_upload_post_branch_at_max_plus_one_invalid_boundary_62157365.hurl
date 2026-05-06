# ── POST /api/upload - branch at max_plus_one_invalid boundary ──
# case_id=TC-62157365
# case_name=POST /api/upload - branch at max_plus_one_invalid boundary
# step_id=step-main
# step_type=test
# technique=boundary_value
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "ILYUfOPfVxuZMfnbVgKKBKcmaHThDumvYBgtnVGhjnPVGeBmGSnwjXFjeojgBxBSehvkPJScHCBTFcjyIabzfzFvTWtmmGsJXlmNIlpLkzqrlyuqKvGoAAOUUwFEBGeoceVrjAMgTmCbeUmYnHVgBpOXAuFUnLPQYGspPdbHIuiUDYqbBJXQtGKAcDLSaGJJLeGIsLZXfWSCbcUflmCylZeRTVGmuNyUFZmpAoeWuylCdFZLpbneeLqzpzLaIKmE",
  "commitSha": "horde",
  "service": "patrol",
  "specContent": "early"
}
```

HTTP 422

