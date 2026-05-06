# ── POST /api/admin/webhooks - name at max_plus_one_invalid boundary ──
# case_id=TC-94214268
# case_name=POST /api/admin/webhooks - name at max_plus_one_invalid boundary
# step_id=step-main
# step_type=test
# technique=boundary_value
# priority=P2

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": [
    "these"
  ],
  "name": "IOYvYIBkAQYqFIqDJMZycrqRFIVCjZIMbSjDHSMaqySSJJGZbEevnwNUYIPXWkWwHWoWMoAdnxnBkAPWCFrpnBgxDdlsucOVjhDdRObECkUodPRyLJNwwstZUaRwXafrnWjLfrJjRGEeTNKnkRrBzcspeyWjjpHjsLvGfcgxXrgoqgfZptELkyLFdklDpBUEtlqfaHPyFoMWMGjhbPWSrFIuUhQHvQOZmItpXjLrWGQNFNXHxaZDTmDNLFhUJSOO",
  "providerType": "infrequently",
  "teamId": "4a6f39f6-5059-431c-b5eb-9711769c6023",
  "url": "http://www.juniorexpedite.com/partnerships"
}
```

HTTP 422

