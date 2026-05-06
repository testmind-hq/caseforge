# ── GET /api/admin/audit-logs - classification tree row 9: [action=user_created] ──
# case_id=TC-e92e324e
# case_name=GET /api/admin/audit-logs - classification tree row 9: [action=user_created]
# step_id=step-main
# step_type=test
# technique=classification_tree
# priority=P2

GET {{base_url}}/api/admin/audit-logs?action=user_created

HTTP 200

[Asserts]
duration < 2000

