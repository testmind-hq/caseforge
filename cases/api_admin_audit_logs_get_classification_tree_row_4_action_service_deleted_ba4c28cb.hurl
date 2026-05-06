# ── GET /api/admin/audit-logs - classification tree row 4: [action=service_deleted] ──
# case_id=TC-ba4c28cb
# case_name=GET /api/admin/audit-logs - classification tree row 4: [action=service_deleted]
# step_id=step-main
# step_type=test
# technique=classification_tree
# priority=P2

GET {{base_url}}/api/admin/audit-logs?action=service_deleted

HTTP 200

[Asserts]
duration < 2000

