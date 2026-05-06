# ── GET /api/admin/audit-logs - classification tree row 10: [action=user_disabled] ──
# case_id=TC-e73ed081
# case_name=GET /api/admin/audit-logs - classification tree row 10: [action=user_disabled]
# step_id=step-main
# step_type=test
# technique=classification_tree
# priority=P2

GET {{base_url}}/api/admin/audit-logs?action=user_disabled

HTTP 200

[Asserts]
duration < 2000

