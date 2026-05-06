# ── GET /api/admin/audit-logs - classification tree row 5: [action=grant_created] ──
# case_id=TC-2874616a
# case_name=GET /api/admin/audit-logs - classification tree row 5: [action=grant_created]
# step_id=step-main
# step_type=test
# technique=classification_tree
# priority=P2

GET {{base_url}}/api/admin/audit-logs?action=grant_created

HTTP 200

[Asserts]
duration < 2000

