# ── [OWASP-API7] DELETE /api/tokens/{id} — injection (sqli) ──
# case_id=TC-e54ea4ce
# case_name=[OWASP-API7] DELETE /api/tokens/{id} — injection (sqli)
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

DELETE {{base_url}}/api/tokens/%27%20OR%201=1--
```json
null
```

HTTP 400

