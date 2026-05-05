// internal/output/render/filename_test.go
package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

func TestFilenameFor(t *testing.T) {
	cases := []struct {
		tc   schema.TestCase
		want string
	}{
		{
			tc: schema.TestCase{
				ID: "TC-3bb73ff9", Title: "POST /users - valid request",
				Steps: []schema.Step{{Method: "POST", Path: "/users"}},
			},
			want: "users_post_valid_request_3bb73ff9",
		},
		{
			tc: schema.TestCase{
				ID: "TC-a1b2c3d4", Title: "GET /pets/{id} - missing id",
				Steps: []schema.Step{{Method: "GET", Path: "/pets/{id}"}},
			},
			want: "pets_id_get_missing_id_a1b2c3d4",
		},
		{
			tc: schema.TestCase{
				ID: "TC-e5f6g7h8", Title: "[owasp] DELETE /admin/teams/{id} - SQL injection",
				Steps: []schema.Step{{Method: "DELETE", Path: "/admin/teams/{id}"}},
			},
			want: "admin_teams_id_delete_owasp_sql_injection_e5f6g7h8",
		},
		{
			tc: schema.TestCase{
				ID: "TC-deadbeef", Title: "PUT /api/admin/users/{id} - boundary value",
				Steps: []schema.Step{{Method: "PUT", Path: "/api/admin/users/{id}"}},
			},
			want: "api_admin_users_id_put_boundary_value_deadbeef",
		},
		{
			// empty title — fall back to just the short ID
			tc:   schema.TestCase{ID: "TC-cafebabe", Title: ""},
			want: "cafebabe",
		},
		{
			tc: schema.TestCase{
				ID: "custom-id", Title: "GET /health - ok",
				Steps: []schema.Step{{Method: "GET", Path: "/health"}},
			},
			want: "health_get_ok_custom-id",
		},
		{
			// chain case: path first, no method, title minus path
			tc: schema.TestCase{
				ID: "TC-abc12345", Title: "CRUD chain: /users", Kind: "chain",
				Steps: []schema.Step{{Method: "POST", Path: "/users"}},
			},
			want: "users_crud_chain_abc12345",
		},
		{
			// no steps — title still parsed for path+method
			tc:   schema.TestCase{ID: "TC-3bb73ff9", Title: "POST /users - valid request"},
			want: "users_post_valid_request_3bb73ff9",
		},
	}

	for _, tc := range cases {
		label := tc.tc.Title
		if label == "" {
			label = "(empty title)"
		}
		t.Run(label, func(t *testing.T) {
			got := FilenameFor(tc.tc)
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

func TestStripMethodPathFromTitle(t *testing.T) {
	cases := []struct {
		title, method, path string
		want                string
	}{
		{
			"POST /users - valid request", "POST", "/users",
			"valid_request",
		},
		{
			"[owasp] DELETE /admin/teams/{id} - SQL injection", "DELETE", "/admin/teams/{id}",
			"owasp_sql_injection",
		},
		{
			"GET /health - ok", "GET", "/health",
			"ok",
		},
		{
			// method+path not found — fallback to full titleSlug
			"some other title", "POST", "/users",
			"some_other_title",
		},
	}
	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			assert.Equal(t, tc.want, stripMethodPathFromTitle(tc.title, tc.method, tc.path))
		})
	}
}

func TestStripPathFromTitle(t *testing.T) {
	cases := []struct {
		title, path string
		want        string
	}{
		{"CRUD chain: /users", "/users", "crud_chain"},
		{"CRUD chain: /store/order", "/store/order", "crud_chain"},
		{"no path here", "/missing", "no_path_here"},
	}
	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			assert.Equal(t, tc.want, stripPathFromTitle(tc.title, tc.path))
		})
	}
}
