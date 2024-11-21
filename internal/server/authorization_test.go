package server_test

// func TestAuthorization(t *testing.T) {
// 	ts, _, err := setupServer()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer ts.Close()
// 	client := ts.Client()
// 	var adminToken string
// 	var nonAdminToken string
// 	t.Run("prepare user accounts and tokens", createAdminAccount(ts.URL, client, &adminToken, &nonAdminToken))

// 	testCases := []struct {
// 		desc   string
// 		method string
// 		path   string
// 		data   any
// 		auth   string
// 		status int
// 	}{
// 		{
// 			desc:   "metrics unreachable without auth",
// 			method: "GET",
// 			path:   "/metrics",
// 			auth:   "",
// 			status: http.StatusOK,
// 		},
// 		{
// 			desc:   "status reachable without auth",
// 			method: "GET",
// 			path:   "/status",
// 			auth:   "",
// 			status: http.StatusOK,
// 		},
// 		{
// 			desc:   "missing endpoints produce 404",
// 			method: "GET",
// 			path:   "/this/path/does/not/exist",
// 			auth:   nonAdminToken,
// 			status: http.StatusNotFound,
// 		},
// 		{
// 			desc:   "nonadmin can't see accounts",
// 			method: "GET",
// 			path:   "/api/v1/employers",
// 			auth:   nonAdminToken,
// 			status: http.StatusForbidden,
// 		},
// 		{
// 			desc:   "admin can see accounts",
// 			method: "GET",
// 			path:   "/api/v1/employers",
// 			auth:   adminToken,
// 			status: http.StatusOK,
// 		},
// 		{
// 			desc:   "nonadmin can't delete admin account",
// 			method: "DELETE",
// 			path:   "/api/v1/employers/1",
// 			auth:   nonAdminToken,
// 			status: http.StatusForbidden,
// 		},
// 		{
// 			desc:   "user can't change admin password",
// 			method: "POST",
// 			path:   "/api/v1/employers/1/change_password",
// 			data:   ChangePasswordRequest{Password: "Pwnd123!"},
// 			auth:   nonAdminToken,
// 			status: http.StatusForbidden,
// 		},
// 		{
// 			desc:   "user can change self password with /me",
// 			method: "POST",
// 			path:   "/api/v1/me/change_password",
// 			data:   ChangePasswordRequest{Password: "BetterPW1!"},
// 			auth:   nonAdminToken,
// 			status: http.StatusOK,
// 		},
// 		{
// 			desc:   "user can login with new password",
// 			method: "POST",
// 			path:   "/api/v1/employers/login",
// 			data: LoginParams{
// 				Email:    "testuser",
// 				Password: "BetterPW1!",
// 			},
// 			auth:   nonAdminToken,
// 			status: http.StatusOK,
// 		},
// 		{
// 			desc:   "user can't list job post",
// 			method: "GET",
// 			path:   "/api/v1/posts",
// 			auth:   nonAdminToken,
// 			status: http.StatusForbidden,
// 		},
// 		{
// 			desc:   "user can't get other user's job post",
// 			method: "GET",
// 			path:   "/api/v1/posts/1",
// 			auth:   nonAdminToken,
// 			status: http.StatusForbidden,
// 		},
// 		{
// 			desc:   "admin can list job post",
// 			method: "GET",
// 			path:   "/api/v1/posts",
// 			auth:   adminToken,
// 			status: http.StatusOK,
// 		},
// 		{
// 			desc:   "admin can't delete itself",
// 			method: "DELETE",
// 			path:   "/api/v1/employers/1",
// 			data:   "",
// 			auth:   adminToken,
// 			status: http.StatusBadRequest,
// 		},
// 		{
// 			desc:   "admin can delete non-admin account",
// 			method: "DELETE",
// 			path:   "/api/v1/employers/2",
// 			auth:   adminToken,
// 			status: http.StatusAccepted,
// 		},
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			body, err := json.Marshal(tC.data)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req, err := http.NewRequest(tC.method, ts.URL+tC.path, bytes.NewReader(body))
// 			req.Header.Add("Authorization", "Bearer "+tC.auth)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			res, err := client.Do(req)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			if res.StatusCode != tC.status {
// 				t.Fatalf("expected status %d, got %d", tC.status, res.StatusCode)
// 			}
// 		})
// 	}
// }
