# ── GET /api/admin/audit-logs - classification tree row 2: [action=spec_uploaded] ──
# case_id=TC-ee7cf268
# case_name=GET /api/admin/audit-logs - classification tree row 2: [action=spec_uploaded]
# step_id=step-main
# step_type=test
# technique=classification_tree
# priority=P2

GET {{base_url}}/api/admin/audit-logs?action=spec_uploaded

HTTP 200

[Asserts]
duration < 2000

