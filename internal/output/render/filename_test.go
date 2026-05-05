// internal/output/render/filename_test.go
package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

func TestFilenameFor(t *testing.T) {
	cases := []struct {
		title string
		id    string
		want  string
	}{
		{
			title: "POST /users - valid request",
			id:    "TC-3bb73ff9",
			want:  "post_users_valid_request_3bb73ff9",
		},
		{
			title: "GET /pets/{id} - missing id",
			id:    "TC-a1b2c3d4",
			want:  "get_pets_id_missing_id_a1b2c3d4",
		},
		{
			title: "[owasp] DELETE /admin/teams/{id} - SQL injection",
			id:    "TC-e5f6g7h8",
			want:  "owasp_delete_admin_teams_id_sql_injection_e5f6g7h8",
		},
		{
			title: "PUT /api/admin/users/{id} - boundary value",
			id:    "TC-deadbeef",
			want:  "put_api_admin_users_id_boundary_value_deadbeef",
		},
		{
			// title is empty — fall back to just the short ID
			title: "",
			id:    "TC-cafebabe",
			want:  "cafebabe",
		},
		{
			// non-TC- prefixed ID passthrough
			title: "GET /health - ok",
			id:    "custom-id",
			want:  "get_health_ok_custom-id",
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			got := FilenameFor(schema.TestCase{ID: tc.id, Title: tc.title})
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestTitleSlug(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"POST /users - valid request", "post_users_valid_request"},
		{"GET /pets/{id}", "get_pets_id"},
		{"[owasp] SQL injection", "owasp_sql_injection"},
		{"", ""},
		{"  leading spaces  ", "leading_spaces"},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			assert.Equal(t, tc.want, titleSlug(tc.in))
		})
	}
}

func TestShortID(t *testing.T) {
	assert.Equal(t, "3bb73ff9", shortID("TC-3bb73ff9"))
	assert.Equal(t, "custom-id", shortID("custom-id"))
	assert.Equal(t, "", shortID("TC-"))
}
