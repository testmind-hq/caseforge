# ── GET /api/admin/audit-logs - classification tree row 6: [action=grant_revoked] ──
# case_id=TC-4511e41f
# case_name=GET /api/admin/audit-logs - classification tree row 6: [action=grant_revoked]
# step_id=step-main
# step_type=test
# technique=classification_tree
# priority=P2

GET {{base_url}}/api/admin/audit-logs?action=grant_revoked

HTTP 200

[Asserts]
duration < 2000

