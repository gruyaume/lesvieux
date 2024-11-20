package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAuthorization(t *testing.T) {
	ts, _, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()
	client := ts.Client()
	var adminToken string
	var nonAdminToken string
	t.Run("prepare user accounts and tokens", prepareUserAccounts(ts.URL, client, &adminToken, &nonAdminToken))

	testCases := []struct {
		desc   string
		method string
		path   string
		data   any
		auth   string
		status int
	}{
		{
			desc:   "metrics unreachable without auth",
			method: "GET",
			path:   "/metrics",
			auth:   "",
			status: http.StatusOK,
		},
		{
			desc:   "status reachable without auth",
			method: "GET",
			path:   "/status",
			auth:   "",
			status: http.StatusOK,
		},
		{
			desc:   "missing endpoints produce 404",
			method: "GET",
			path:   "/this/path/does/not/exist",
			auth:   nonAdminToken,
			status: http.StatusNotFound,
		},
		{
			desc:   "nonadmin can't see accounts",
			method: "GET",
			path:   "/api/v1/accounts",
			auth:   nonAdminToken,
			status: http.StatusForbidden,
		},
		{
			desc:   "admin can see accounts",
			method: "GET",
			path:   "/api/v1/accounts",
			auth:   adminToken,
			status: http.StatusOK,
		},
		{
			desc:   "nonadmin can't delete admin account",
			method: "DELETE",
			path:   "/api/v1/accounts/1",
			auth:   nonAdminToken,
			status: http.StatusForbidden,
		},
		{
			desc:   "user can't change admin password",
			method: "POST",
			path:   "/api/v1/accounts/1/change_password",
			data:   ChangePasswordRequest{Password: "Pwnd123!"},
			auth:   nonAdminToken,
			status: http.StatusForbidden,
		},
		{
			desc:   "user can change self password with /me",
			method: "POST",
			path:   "/api/v1/me/change_password",
			data:   ChangePasswordRequest{Password: "BetterPW1!"},
			auth:   nonAdminToken,
			status: http.StatusOK,
		},
		{
			desc:   "user can login with new password",
			method: "POST",
			path:   "/api/v1/login",
			data: LoginParams{
				Username: "testuser",
				Password: "BetterPW1!",
			},
			auth:   nonAdminToken,
			status: http.StatusOK,
		},
		{
			desc:   "admin can create own blog post",
			method: "POST",
			path:   "/api/v1/me/posts",
			auth:   adminToken,
			status: http.StatusCreated,
		},
		{
			desc:   "admin can edit own blog post",
			method: "PUT",
			path:   "/api/v1/me/posts/1",
			data: UpdateBlogPostParams{
				Title:   "My Title",
				Content: "My Content",
				Status:  "published",
			},
			auth:   adminToken,
			status: http.StatusOK,
		},
		{
			desc:   "admin can list own blog posts",
			method: "GET",
			path:   "/api/v1/me/posts",
			auth:   adminToken,
			status: http.StatusOK,
		},
		{
			desc:   "admin can get own blog post",
			method: "GET",
			path:   "/api/v1/me/posts/1",
			auth:   adminToken,
			status: http.StatusOK,
		},
		{
			desc:   "user can't list blog post",
			method: "GET",
			path:   "/api/v1/posts",
			auth:   nonAdminToken,
			status: http.StatusForbidden,
		},
		{
			desc:   "user can't get other user's blog post",
			method: "GET",
			path:   "/api/v1/posts/1",
			auth:   nonAdminToken,
			status: http.StatusForbidden,
		},
		{
			desc:   "user can't get other user's blog post using me",
			method: "GET",
			path:   "/api/v1/me/posts/1",
			auth:   nonAdminToken,
			status: http.StatusForbidden,
		},
		{
			desc:   "user can't edit other user's blog post using me",
			method: "PUT",
			path:   "/api/v1/me/posts/1",
			auth:   nonAdminToken,
			status: http.StatusForbidden,
		},
		{
			desc:   "user can create blog post",
			method: "POST",
			path:   "/api/v1/me/posts",
			auth:   nonAdminToken,
			status: http.StatusCreated,
		},
		{
			desc:   "user can edit own blog post",
			method: "PUT",
			path:   "/api/v1/me/posts/2",
			data: UpdateBlogPostParams{
				Title:   "My Title",
				Content: "My Content",
				Status:  "draft",
			},
			auth:   nonAdminToken,
			status: http.StatusOK,
		},
		{
			desc:   "user can list own blog posts",
			method: "GET",
			path:   "/api/v1/me/posts",
			auth:   nonAdminToken,
			status: http.StatusOK,
		},
		{
			desc:   "user can get own blog post",
			method: "GET",
			path:   "/api/v1/me/posts/2",
			auth:   nonAdminToken,
			status: http.StatusOK,
		},
		{
			desc:   "admin can list blog post",
			method: "GET",
			path:   "/api/v1/posts",
			auth:   adminToken,
			status: http.StatusOK,
		},
		{
			desc:   "admin can get other user's blog post",
			method: "GET",
			path:   "/api/v1/posts/2",
			auth:   adminToken,
			status: http.StatusOK,
		},
		{
			desc:   "admin can't delete itself",
			method: "DELETE",
			path:   "/api/v1/accounts/1",
			data:   "",
			auth:   adminToken,
			status: http.StatusBadRequest,
		},
		{
			desc:   "admin can delete non-admin account",
			method: "DELETE",
			path:   "/api/v1/accounts/2",
			auth:   adminToken,
			status: http.StatusAccepted,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			body, err := json.Marshal(tC.data)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(tC.method, ts.URL+tC.path, bytes.NewReader(body))
			req.Header.Add("Authorization", "Bearer "+tC.auth)
			if err != nil {
				t.Fatal(err)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			if res.StatusCode != tC.status {
				t.Fatalf("expected status %d, got %d", tC.status, res.StatusCode)
			}
		})
	}
}
