# ── POST /api/upload - specContent at max_valid boundary ──
# case_id=TC-201ba23b
# case_name=POST /api/upload - specContent at max_valid boundary
# step_id=step-main
# step_type=test
# technique=boundary_value
# priority=P1

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "annoying",
  "commitSha": "horde",
  "service": "patrol",
  "specContent": "MvxueBBOuEUznvCnujHEfhfJEmIkMiFxMUaMDQYopjbpdETOJXbhaSibxhItFKowWSgvVTsEKoRBvRboGZCrpNFYbErOCedxMcVAnLzDekWtkEvgLpSZAGaDLsFRvNWihavpvGqXfpluZjqXgXkvQZEpaaHgrFeEHQhhHsZqkGppwxBdpFmjShygsygoqyopydhyLxSwTwouvqLXCFkgNFkmEiZKFOzPodlBbQdZyQXKtqOjjyxMqTwcyXFgxoI"
}
```

HTTP 200

[Asserts]
duration < 2000

