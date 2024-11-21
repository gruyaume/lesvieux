package server_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

type CreateAdminAccountParams struct {
	Id       int    `json:"id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangeAdminPasswordRequest struct {
	Password string `json:"password"`
}

type GetAdminAccountResponseResult struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type GetAdminAccountResponse struct {
	Result GetAdminAccountResponseResult `json:"result"`
	Error  string                        `json:"error,omitempty"`
}

type CreateAdminAccountResponseResult struct {
	Id int `json:"id"`
}

type CreateAdminAccountResponse struct {
	Result CreateAdminAccountResponseResult `json:"result"`
	Error  string                           `json:"error,omitempty"`
}

type AdminPasswordResponseResult struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

type ChangeAdminPasswordResponse struct {
	Error  string                      `json:"error"`
	Result AdminPasswordResponseResult `json:"result"`
}

type DeleteAdminAccountResponseResult struct {
	Id int `json:"id"`
}

type DeleteAdminAccountResponse struct {
	Error  string                           `json:"error"`
	Result DeleteAdminAccountResponseResult `json:"result"`
}

func createAdminAccount(url string, client *http.Client, adminToken string, data *CreateAdminAccountParams) (int, *CreateAdminAccountResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest("POST", url+"/api/v1/admin/accounts", strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+adminToken)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var createResponse CreateAdminAccountResponse
	if err := json.NewDecoder(res.Body).Decode(&createResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &createResponse, nil
}

func getAdminAccount(url string, client *http.Client, token string, id string) (int, *GetAdminAccountResponse, error) {
	req, err := http.NewRequest("GET", url+"/api/v1/admin/accounts/"+id, nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var getResponse GetAdminAccountResponse
	if err := json.NewDecoder(res.Body).Decode(&getResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &getResponse, nil
}

func getMyAdminAccount(url string, client *http.Client, token string) (int, *GetAdminAccountResponse, error) {
	req, err := http.NewRequest("GET", url+"/api/v1/admin/accounts/me", nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var getResponse GetAdminAccountResponse
	if err := json.NewDecoder(res.Body).Decode(&getResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &getResponse, nil
}

func deleteAdminAccount(url string, client *http.Client, token string, id string) (int, *DeleteAdminAccountResponse, error) {
	req, err := http.NewRequest("DELETE", url+"/api/v1/admin/accounts/"+id, nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var deleteResponse DeleteAdminAccountResponse
	if err := json.NewDecoder(res.Body).Decode(&deleteResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &deleteResponse, nil
}

func changeAdminAccountPassword(url string, client *http.Client, token string, id string, data *ChangeAdminPasswordRequest) (int, *ChangeAdminPasswordResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest("POST", url+"/api/v1/admin/accounts/"+id+"/change_password", strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var changeResponse ChangeAdminPasswordResponse
	if err := json.NewDecoder(res.Body).Decode(&changeResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &changeResponse, nil
}

func changeMyAdminAccountPassword(url string, client *http.Client, token string, data *ChangeAdminPasswordRequest) (int, *ChangeAdminPasswordResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest("POST", url+"/api/v1/admin/accounts/me/change_password", strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var changeResponse ChangeAdminPasswordResponse
	if err := json.NewDecoder(res.Body).Decode(&changeResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &changeResponse, nil
}

func TestUsersHandlersCreateAccount(t *testing.T) {
	ts, _, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()
	client := ts.Client()
	var token string
	t.Run("prepare user accounts and tokens", prepareAdminAccount(ts.URL, client, &token))

	testCases := []struct {
		desc             string
		auth             string
		data             CreateAdminAccountParams
		expectedResponse CreateAdminAccountResponse
		status           int
	}{
		{
			desc: "Admin create user - success",
			data: CreateAdminAccountParams{Email: "testuser@guillaume.com", Password: "Password1!"},
			auth: token,
			expectedResponse: CreateAdminAccountResponse{
				Result: CreateAdminAccountResponseResult{
					Id: 2,
				},
			},
			status: http.StatusCreated,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := createAdminAccount(ts.URL, client, tC.auth, &tC.data)
			if err != nil {
				t.Fatalf("couldn't create account: %s", err)
			}
			if statusCode != tC.status {
				t.Fatalf("expected status %d, got %d", tC.status, statusCode)
			}
			if resp.Error != tC.expectedResponse.Error {
				t.Fatalf("expected error %q, got %q", tC.expectedResponse.Error, resp.Error)
			}
			if resp.Result != tC.expectedResponse.Result {
				t.Fatalf("expected result %v, got %v", tC.expectedResponse.Result, resp.Result)
			}
		})
	}
}

func TestHandlersGetAdminAccount(t *testing.T) {
	ts, _, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()
	client := ts.Client()
	var token string
	t.Run("prepare user accounts and tokens", prepareAdminAccount(ts.URL, client, &token))

	testCases := []struct {
		desc             string
		id               string
		auth             string
		expectedResponse GetAdminAccountResponse
		status           int
	}{
		{
			desc: "Admin get admin user - success",
			id:   "1",
			auth: token,
			expectedResponse: GetAdminAccountResponse{
				Result: GetAdminAccountResponseResult{
					Id: 1, Email: "testadmin",
				},
			},
			status: http.StatusOK,
		},
		{
			desc:             "Admin get inexistent user - fail",
			id:               "300",
			auth:             token,
			expectedResponse: GetAdminAccountResponse{Error: "Admin Account not found"},
			status:           http.StatusNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := getAdminAccount(ts.URL, client, tC.auth, tC.id)
			if err != nil {
				t.Fatalf("couldn't get account: %s", err)
			}
			if statusCode != tC.status {
				t.Fatalf("expected status %d, got %d", tC.status, statusCode)
			}
			if resp.Error != tC.expectedResponse.Error {
				t.Fatalf("expected error %q, got %q", tC.expectedResponse.Error, resp.Error)
			}
			if resp.Result != tC.expectedResponse.Result {
				t.Fatalf("expected result %v, got %v", tC.expectedResponse.Result, resp.Result)
			}
		})
	}
}

func TestHandlersGetMyAdminAccount(t *testing.T) {
	ts, _, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()
	client := ts.Client()
	var token string
	t.Run("prepare user accounts and tokens", prepareAdminAccount(ts.URL, client, &token))

	testCases := []struct {
		desc             string
		auth             string
		expectedResponse GetAdminAccountResponse
		status           int
	}{
		{
			desc: "Get own user - success",
			auth: token,
			expectedResponse: GetAdminAccountResponse{
				Result: GetAdminAccountResponseResult{
					Id: 1, Email: "testadmin",
				},
			},
			status: http.StatusOK,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := getMyAdminAccount(ts.URL, client, tC.auth)
			if err != nil {
				t.Fatalf("couldn't get account: %s", err)
			}
			if statusCode != tC.status {
				t.Fatalf("expected status %d, got %d", tC.status, statusCode)
			}
			if resp.Error != tC.expectedResponse.Error {
				t.Fatalf("expected error %q, got %q", tC.expectedResponse.Error, resp.Error)
			}
			if resp.Result != tC.expectedResponse.Result {
				t.Fatalf("expected result %v, got %v", tC.expectedResponse.Result, resp.Result)
			}
		})
	}
}

func TestUsersHandlersDeleteAccount(t *testing.T) {
	ts, _, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()
	client := ts.Client()
	var token string
	t.Run("prepare user accounts and tokens", prepareAdminAccount(ts.URL, client, &token))

	testCases := []struct {
		desc             string
		method           string
		id               string
		auth             string
		expectedResponse DeleteAdminAccountResponse
		status           int
	}{

		{
			desc:             "Admin delete non-existent user - failure",
			id:               "123",
			auth:             token,
			expectedResponse: DeleteAdminAccountResponse{Error: "Admin Account not found"},
			status:           http.StatusNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := deleteAdminAccount(ts.URL, client, tC.auth, tC.id)
			if err != nil {
				t.Fatalf("couldn't delete account: %s", err)
			}
			if statusCode != tC.status {
				t.Fatalf("expected status %d, got %d", tC.status, statusCode)
			}

			if resp.Error != tC.expectedResponse.Error {
				t.Fatalf("expected error %q, got %q", tC.expectedResponse.Error, resp.Error)
			}
			if resp.Result != tC.expectedResponse.Result {
				t.Fatalf("expected result %v, got %v", tC.expectedResponse.Result, resp.Result)
			}
		})
	}
}

func TestHandlersChangeAdminAccountPassword(t *testing.T) {
	ts, _, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()
	client := ts.Client()
	var token string
	t.Run("prepare user accounts and tokens", prepareAdminAccount(ts.URL, client, &token))

	testCases := []struct {
		desc             string
		id               string
		data             ChangeAdminPasswordRequest
		auth             string
		expectedResponse ChangeAdminPasswordResponse
		status           int
	}{
		{
			desc:             "Admin change password with id - success",
			id:               "1",
			data:             ChangeAdminPasswordRequest{Password: "newPassword1"},
			auth:             token,
			expectedResponse: ChangeAdminPasswordResponse{Result: AdminPasswordResponseResult{Id: 1}},
			status:           http.StatusOK,
		},
		{
			desc:             "Admin change non-existing user password - failure",
			id:               "100",
			data:             ChangeAdminPasswordRequest{Password: "newPassword1"},
			auth:             token,
			expectedResponse: ChangeAdminPasswordResponse{Error: "Admin Account not found"},
			status:           http.StatusNotFound,
		},
		{
			desc:             "Admin change password without password - failure",
			id:               "1",
			data:             ChangeAdminPasswordRequest{},
			auth:             token,
			expectedResponse: ChangeAdminPasswordResponse{Error: "Password is required"},
			status:           http.StatusBadRequest,
		},
		{
			desc:             "Admin change password with bad password - failure",
			id:               "1",
			data:             ChangeAdminPasswordRequest{Password: "password"},
			auth:             token,
			expectedResponse: ChangeAdminPasswordResponse{Error: "Password must have 8 or more characters, must include at least one capital letter, one lowercase letter, and either a number or a symbol."},
			status:           http.StatusBadRequest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := changeAdminAccountPassword(ts.URL, client, tC.auth, tC.id, &tC.data)
			if err != nil {
				t.Fatalf("couldn't change password: %s", err)
			}

			if statusCode != tC.status {
				t.Fatalf("expected status %d, got %d", tC.status, statusCode)
			}

			if resp.Error != tC.expectedResponse.Error {
				t.Fatalf("expected error %q, got %q", tC.expectedResponse.Error, resp.Error)
			}
			if resp.Result != tC.expectedResponse.Result {
				t.Fatalf("expected result %v, got %v", tC.expectedResponse.Result, resp.Result)
			}
		})
	}
}

func TestUsersHandlersChangeMyAdminPassword(t *testing.T) {
	ts, _, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()
	client := ts.Client()
	var token string
	t.Run("prepare user accounts and tokens", prepareAdminAccount(ts.URL, client, &token))

	testCases := []struct {
		desc             string
		data             ChangeAdminPasswordRequest
		auth             string
		expectedResponse ChangeAdminPasswordResponse
		status           int
	}{
		{
			desc:             "Admin change password with me - success",
			data:             ChangeAdminPasswordRequest{Password: "newPassword1"},
			auth:             token,
			expectedResponse: ChangeAdminPasswordResponse{Result: AdminPasswordResponseResult{Id: 1}},
			status:           http.StatusOK,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := changeMyAdminAccountPassword(ts.URL, client, tC.auth, &tC.data)
			if err != nil {
				t.Fatalf("couldn't change password: %s", err)
			}

			if statusCode != tC.status {
				t.Fatalf("expected status %d, got %d", tC.status, statusCode)
			}

			if resp.Error != tC.expectedResponse.Error {
				t.Fatalf("expected error %q, got %q", tC.expectedResponse.Error, resp.Error)
			}
			if resp.Result != tC.expectedResponse.Result {
				t.Fatalf("expected result %v, got %v", tC.expectedResponse.Result, resp.Result)
			}
		})
	}
}
