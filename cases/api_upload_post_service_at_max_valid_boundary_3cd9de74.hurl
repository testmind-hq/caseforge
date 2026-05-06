# ── POST /api/upload - service at max_valid boundary ──
# case_id=TC-3cd9de74
# case_name=POST /api/upload - service at max_valid boundary
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
  "service": "atLOmtVVmlQhFvFrwuMTJjhgqzDQgMAKdxkeUnYswKYRxCFECDdRtuhENDYOeachFgpnTjKElKhbRGMNBMqtQcJeLmJEdXosWDnsTCROKgowmZMFmjZPjXeSVkrLtqyrTdhcTIoNWdfwRXnmvZQoROrQlafSbnQScDRKBvbCIsqPEGzseScyClXaqHCuhtwbNgwbAjmxZkPvBMGOxVbdVVDWFWdnUugVnZaDTXdkaRzAOYonKbCYZPlwlDZDKdT",
  "specContent": "early"
}
```

HTTP 200

[Asserts]
duration < 2000

