# ── GET /api/admin/audit-logs - classification tree row 8: [action=token_revoked] ──
# case_id=TC-5a6e9137
# case_name=GET /api/admin/audit-logs - classification tree row 8: [action=token_revoked]
# step_id=step-main
# step_type=test
# technique=classification_tree
# priority=P2

GET {{base_url}}/api/admin/audit-logs?action=token_revoked

HTTP 200

[Asserts]
duration < 2000

