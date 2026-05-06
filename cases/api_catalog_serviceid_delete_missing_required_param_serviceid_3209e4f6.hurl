# ── DELETE /api/catalog/:serviceId - missing required param "serviceId" ──
# case_id=TC-3209e4f6
# case_name=DELETE /api/catalog/:serviceId - missing required param "serviceId"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

DELETE {{base_url}}/api/catalog/:serviceId

HTTP 422

