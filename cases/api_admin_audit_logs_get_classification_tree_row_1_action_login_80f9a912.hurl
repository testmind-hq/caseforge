# ── GET /api/admin/audit-logs - classification tree row 1: [action=login] ──
# case_id=TC-80f9a912
# case_name=GET /api/admin/audit-logs - classification tree row 1: [action=login]
# step_id=step-main
# step_type=test
# technique=classification_tree
# priority=P2

GET {{base_url}}/api/admin/audit-logs?action=login

HTTP 200

[Asserts]
duration < 2000

