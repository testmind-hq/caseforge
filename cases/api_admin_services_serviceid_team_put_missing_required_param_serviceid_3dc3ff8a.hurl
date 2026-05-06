# ── PUT /api/admin/services/{serviceId}/team - missing required param "serviceId" ──
# case_id=TC-3dc3ff8a
# case_name=PUT /api/admin/services/{serviceId}/team - missing required param "serviceId"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

PUT {{base_url}}/api/admin/services/1/team

HTTP 422

