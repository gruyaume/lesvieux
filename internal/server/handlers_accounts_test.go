package server_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

type CreateAccountParams struct {
	Id       int    `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	Password string `json:"password"`
}

type GetAccountResponseResult struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
}

type GetAccountResponse struct {
	Result GetAccountResponseResult `json:"result"`
	Error  string                   `json:"error,omitempty"`
}

type CreateAccountResponseResult struct {
	Id int `json:"id"`
}

type CreateAccountResponse struct {
	Result CreateAccountResponseResult `json:"result"`
	Error  string                      `json:"error,omitempty"`
}

type PasswordResponseResult struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

type ChangePasswordResponse struct {
	Error  string                 `json:"error"`
	Result PasswordResponseResult `json:"result"`
}

type DeleteAccountResponseResult struct {
	Id int `json:"id"`
}

type DeleteAccountResponse struct {
	Error  string                      `json:"error"`
	Result DeleteAccountResponseResult `json:"result"`
}

func createAccount(url string, client *http.Client, adminToken string, data *CreateAccountParams) (int, *CreateAccountResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest("POST", url+"/api/v1/accounts", strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+adminToken)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var createResponse CreateAccountResponse
	if err := json.NewDecoder(res.Body).Decode(&createResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &createResponse, nil
}

func getAccount(url string, client *http.Client, token string, id string) (int, *GetAccountResponse, error) {
	req, err := http.NewRequest("GET", url+"/api/v1/accounts/"+id, nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var getResponse GetAccountResponse
	if err := json.NewDecoder(res.Body).Decode(&getResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &getResponse, nil
}

func getMyAccount(url string, client *http.Client, token string) (int, *GetAccountResponse, error) {
	req, err := http.NewRequest("GET", url+"/api/v1/me", nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var getResponse GetAccountResponse
	if err := json.NewDecoder(res.Body).Decode(&getResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &getResponse, nil
}

func deleteAccount(url string, client *http.Client, token string, id string) (int, *DeleteAccountResponse, error) {
	req, err := http.NewRequest("DELETE", url+"/api/v1/accounts/"+id, nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var deleteResponse DeleteAccountResponse
	if err := json.NewDecoder(res.Body).Decode(&deleteResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &deleteResponse, nil
}

func changeAccountPassword(url string, client *http.Client, token string, id string, data *ChangePasswordRequest) (int, *ChangePasswordResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest("POST", url+"/api/v1/accounts/"+id+"/change_password", strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var changeResponse ChangePasswordResponse
	if err := json.NewDecoder(res.Body).Decode(&changeResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &changeResponse, nil
}

func changeMyAccountPassword(url string, client *http.Client, token string, data *ChangePasswordRequest) (int, *ChangePasswordResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest("POST", url+"/api/v1/me/change_password", strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var changeResponse ChangePasswordResponse
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
	var adminToken string
	var nonAdminToken string
	t.Run("prepare user accounts and tokens", prepareUserAccounts(ts.URL, client, &adminToken, &nonAdminToken))

	testCases := []struct {
		desc             string
		auth             string
		data             CreateAccountParams
		expectedResponse CreateAccountResponse
		status           int
	}{
		{
			desc: "Admin create user - success",
			data: CreateAccountParams{Username: "testuser2", Password: "Password1!"},
			auth: adminToken,
			expectedResponse: CreateAccountResponse{
				Result: CreateAccountResponseResult{
					Id: 3,
				},
			},
			status: http.StatusCreated,
		},
		{
			desc:             "Non-Admin create user - fail",
			data:             CreateAccountParams{Username: "testuser3", Password: "Password1!"},
			auth:             nonAdminToken,
			expectedResponse: CreateAccountResponse{Error: "forbidden: admin access required"},
			status:           http.StatusForbidden,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := createAccount(ts.URL, client, tC.auth, &tC.data)
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

func TestUsersHandlersGetAccount(t *testing.T) {
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
		desc             string
		id               string
		auth             string
		expectedResponse GetAccountResponse
		status           int
	}{
		{
			desc: "Admin get admin user - success",
			id:   "1",
			auth: adminToken,
			expectedResponse: GetAccountResponse{
				Result: GetAccountResponseResult{
					Id: 1, Username: "testadmin", Role: 1,
				},
			},
			status: http.StatusOK,
		},

		{
			desc: "Admin get normal user success",
			id:   "2",
			auth: adminToken,
			expectedResponse: GetAccountResponse{
				Result: GetAccountResponseResult{
					Id:       2,
					Username: "testuser",
					Role:     0,
				},
			},
			status: http.StatusOK,
		},
		{
			desc:             "Admin get inexistent user - fail",
			id:               "300",
			auth:             adminToken,
			expectedResponse: GetAccountResponse{Error: "Account not found"},
			status:           http.StatusNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := getAccount(ts.URL, client, tC.auth, tC.id)
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

func TestUsersHandlersGetMyAccount(t *testing.T) {
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
		desc             string
		auth             string
		expectedResponse GetAccountResponse
		status           int
	}{
		{
			desc: "Non-Admin get own user - success",
			auth: nonAdminToken,
			expectedResponse: GetAccountResponse{
				Result: GetAccountResponseResult{
					Id: 2, Username: "testuser", Role: 0,
				},
			},
			status: http.StatusOK,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := getMyAccount(ts.URL, client, tC.auth)
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
	var adminToken string
	var nonAdminToken string
	t.Run("prepare user accounts and tokens", prepareUserAccounts(ts.URL, client, &adminToken, &nonAdminToken))

	testCases := []struct {
		desc             string
		method           string
		id               string
		auth             string
		expectedResponse DeleteAccountResponse
		status           int
	}{
		{
			desc:             "Non-admin delete user - failure",
			id:               "2",
			auth:             nonAdminToken,
			expectedResponse: DeleteAccountResponse{Error: "forbidden: admin access required"},
			status:           http.StatusForbidden,
		},
		{
			desc:             "Admin delete user - success",
			id:               "2",
			auth:             adminToken,
			expectedResponse: DeleteAccountResponse{Result: DeleteAccountResponseResult{Id: 2}},
			status:           http.StatusAccepted,
		},
		{
			desc:             "Admin delete non-existent user - failure",
			id:               "123",
			auth:             adminToken,
			expectedResponse: DeleteAccountResponse{Error: "Account not found"},
			status:           http.StatusNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := deleteAccount(ts.URL, client, tC.auth, tC.id)
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

func TestUsersHandlersChangePassword(t *testing.T) {
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
		desc             string
		id               string
		data             ChangePasswordRequest
		auth             string
		expectedResponse ChangePasswordResponse
		status           int
	}{
		{
			desc:             "Admin change password with id - success",
			id:               "1",
			data:             ChangePasswordRequest{Password: "newPassword1"},
			auth:             adminToken,
			expectedResponse: ChangePasswordResponse{Result: PasswordResponseResult{Id: 1}},
			status:           http.StatusOK,
		},
		{
			desc:             "Non admin change other user password - failure",
			id:               "1",
			data:             ChangePasswordRequest{Password: "newPassword1"},
			auth:             nonAdminToken,
			expectedResponse: ChangePasswordResponse{Error: "forbidden: admin access required"},
			status:           http.StatusForbidden,
		},
		{
			desc:             "Admin change non-existing user password - failure",
			id:               "100",
			data:             ChangePasswordRequest{Password: "newPassword1"},
			auth:             adminToken,
			expectedResponse: ChangePasswordResponse{Error: "Account not found"},
			status:           http.StatusNotFound,
		},
		{
			desc:             "Admin change password without password - failure",
			id:               "1",
			data:             ChangePasswordRequest{},
			auth:             adminToken,
			expectedResponse: ChangePasswordResponse{Error: "Password is required"},
			status:           http.StatusBadRequest,
		},
		{
			desc:             "Admin change password with bad password - failure",
			id:               "1",
			data:             ChangePasswordRequest{Password: "password"},
			auth:             adminToken,
			expectedResponse: ChangePasswordResponse{Error: "Password must have 8 or more characters, must include at least one capital letter, one lowercase letter, and either a number or a symbol."},
			status:           http.StatusBadRequest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := changeAccountPassword(ts.URL, client, tC.auth, tC.id, &tC.data)
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

func TestUsersHandlersChangeMyPassword(t *testing.T) {
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
		desc             string
		data             ChangePasswordRequest
		auth             string
		expectedResponse ChangePasswordResponse
		status           int
	}{
		{
			desc:             "Admin change password with me - success",
			data:             ChangePasswordRequest{Password: "newPassword1"},
			auth:             adminToken,
			expectedResponse: ChangePasswordResponse{Result: PasswordResponseResult{Id: 1}},
			status:           http.StatusOK,
		},
		{
			desc:             "Non admin change own password - success",
			data:             ChangePasswordRequest{Password: "newPassword1"},
			auth:             nonAdminToken,
			expectedResponse: ChangePasswordResponse{Result: PasswordResponseResult{Id: 2}},
			status:           http.StatusOK,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := changeMyAccountPassword(ts.URL, client, tC.auth, &tC.data)
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
