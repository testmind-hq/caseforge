# ── GET /api/admin/audit-logs - classification tree row 7: [action=token_created] ──
# case_id=TC-e290ff04
# case_name=GET /api/admin/audit-logs - classification tree row 7: [action=token_created]
# step_id=step-main
# step_type=test
# technique=classification_tree
# priority=P2

GET {{base_url}}/api/admin/audit-logs?action=token_created

HTTP 200

[Asserts]
duration < 2000

