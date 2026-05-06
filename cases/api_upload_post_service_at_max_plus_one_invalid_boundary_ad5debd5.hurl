# ── POST /api/upload - service at max_plus_one_invalid boundary ──
# case_id=TC-ad5debd5
# case_name=POST /api/upload - service at max_plus_one_invalid boundary
# step_id=step-main
# step_type=test
# technique=boundary_value
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "annoying",
  "commitSha": "horde",
  "service": "UqQKQdxIBaEEFIOlbucPEjkejpJhtGCnYytkTfHBnTHmoeamHxyFTtNkqceSxPhYjEZfVjxnkUrCXnzCRdtVbcomgJaqcHidTZbQHOJgFusDCcCXqQuHRTajulzyqxxOFgJZTIrWbrgvHDgjlzyuuBztsMwepFaVmllpLTRwhONiNNZZDMtJFSySHEyRBmGBvFwEkoyGZJSFbcrJaJVmftRoXuHFuUwcKLaJFIIGOYYgsNiAMNTBUcmdjtEEKcrT",
  "specContent": "early"
}
```

HTTP 422

