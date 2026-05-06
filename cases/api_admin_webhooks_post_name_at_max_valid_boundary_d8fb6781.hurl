# ── POST /api/admin/webhooks - name at max_valid boundary ──
# case_id=TC-d8fb6781
# case_name=POST /api/admin/webhooks - name at max_valid boundary
# step_id=step-main
# step_type=test
# technique=boundary_value
# priority=P1

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": [
    "these"
  ],
  "name": "SncWFCUvZpQFNFdrRgNJvYbFANxRmLnQRwBDZqHrTHNxToOSzvIyMmzYXYNlTmqxqecveYPPJkHsbPGoaolHtERzLSSWSCxHgCRyXtiMrbXGLHWZPsGbytTNsOuzeJeHwrLudLzbVBdbBDdVDJAEXLewLKAlJsnbYaiuzbPulctRaehbdWqhpaxcUFmpSCgDEsQEUPqkVaYFLwaCaeKPlKLmHypHEUNlnmuYwzseXfFSYIVfMKOFtwTgnGGRbhK",
  "providerType": "infrequently",
  "teamId": "4a6f39f6-5059-431c-b5eb-9711769c6023",
  "url": "http://www.juniorexpedite.com/partnerships"
}
```

HTTP 200

[Asserts]
duration < 2000

