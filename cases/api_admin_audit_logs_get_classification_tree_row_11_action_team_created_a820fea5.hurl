# ── GET /api/admin/audit-logs - classification tree row 11: [action=team_created] ──
# case_id=TC-a820fea5
# case_name=GET /api/admin/audit-logs - classification tree row 11: [action=team_created]
# step_id=step-main
# step_type=test
# technique=classification_tree
# priority=P2

GET {{base_url}}/api/admin/audit-logs?action=team_created

HTTP 200

[Asserts]
duration < 2000

