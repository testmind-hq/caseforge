# ── DELETE /api/catalog/:serviceId - IDOR serviceId=00000000-0000-0000-0000-000000000001 (alt_uuid) ──
# case_id=TC-e72a9984
# case_name=DELETE /api/catalog/:serviceId - IDOR serviceId=00000000-0000-0000-0000-000000000001 (alt_uuid)
# step_id=step-main
# step_type=test
# technique=idor
# priority=P1

DELETE {{base_url}}/api/catalog/:serviceId

HTTP *

[Asserts]
status >= 400
status < 500

