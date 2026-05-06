# ── GET /api/admin/audit-logs - classification tree row 3: [action=spec_updated] ──
# case_id=TC-df4697d4
# case_name=GET /api/admin/audit-logs - classification tree row 3: [action=spec_updated]
# step_id=step-main
# step_type=test
# technique=classification_tree
# priority=P2

GET {{base_url}}/api/admin/audit-logs?action=spec_updated

HTTP 200

[Asserts]
duration < 2000

