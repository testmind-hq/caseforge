# ── POST /api/tokens - name at max_valid boundary ──
# case_id=TC-a0247f03
# case_name=POST /api/tokens - name at max_valid boundary
# step_id=step-main
# step_type=test
# technique=boundary_value
# priority=P1

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "dIcVzeAXIpwOMzbhuWAKvYpdHpXhDnlquznBMpHNObsplNJMCmfagUMlgmyfFcxjiOSjnDPJMExECRCIPMONUmxCjiZwOKphjBRzxRgqBHCPWiUvPVxGpuIuOwqcjGDtPEXvUFwTFgNBEKmwQejgeRCcxYCgaGRusgCHYhGuMkhuWBKpkpOWZMOWQrWAqMGwVOnWXHenTnRwxoXQNWVzoLuAeLfEUWmvtOaUOzDopkvdpjDJgEGrzToimadBCbq",
  "scope": "read"
}
```

HTTP 200

[Asserts]
duration < 2000

