package server_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

type CreateEmployerAccountParams struct {
	Id       int    `json:"id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangeEmployerPasswordRequest struct {
	Password string `json:"password"`
}

type GetEmployerAccountResponseResult struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type GetEmployerAccountResponse struct {
	Result GetEmployerAccountResponseResult `json:"result"`
	Error  string                           `json:"error,omitempty"`
}

type CreateEmployerAccountResponseResult struct {
	Id int `json:"id"`
}

type CreateEmployerAccountResponse struct {
	Result CreateEmployerAccountResponseResult `json:"result"`
	Error  string                              `json:"error,omitempty"`
}

type EmployerPasswordResponseResult struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

type ChangeEmployerPasswordResponse struct {
	Error  string                         `json:"error"`
	Result EmployerPasswordResponseResult `json:"result"`
}

type DeleteEmployerAccountResponseResult struct {
	Id int `json:"id"`
}

type DeleteEmployerAccountResponse struct {
	Error  string                              `json:"error"`
	Result DeleteEmployerAccountResponseResult `json:"result"`
}

func createEmployerAccount(url string, client *http.Client, adminToken string, employerID string, data *CreateEmployerAccountParams) (int, *CreateEmployerAccountResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest("POST", url+"/api/v1/employers/"+employerID+"/accounts", strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+adminToken)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var createResponse CreateEmployerAccountResponse
	if err := json.NewDecoder(res.Body).Decode(&createResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &createResponse, nil
}

func getEmployerAccount(url string, client *http.Client, token string, id string) (int, *GetEmployerAccountResponse, error) {
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
	var getResponse GetEmployerAccountResponse
	if err := json.NewDecoder(res.Body).Decode(&getResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &getResponse, nil
}

func getMyEmployerAccount(url string, client *http.Client, token string) (int, *GetEmployerAccountResponse, error) {
	req, err := http.NewRequest("GET", url+"/api/v1/employers/accounts/me", nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var getResponse GetEmployerAccountResponse
	if err := json.NewDecoder(res.Body).Decode(&getResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &getResponse, nil
}

func deleteEmployerAccount(url string, client *http.Client, token string, id string) (int, *DeleteEmployerAccountResponse, error) {
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
	var deleteResponse DeleteEmployerAccountResponse
	if err := json.NewDecoder(res.Body).Decode(&deleteResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &deleteResponse, nil
}

func changeEmployerAccountPassword(url string, client *http.Client, token string, id string, data *ChangeEmployerPasswordRequest) (int, *ChangeEmployerPasswordResponse, error) {
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
	var changeResponse ChangeEmployerPasswordResponse
	if err := json.NewDecoder(res.Body).Decode(&changeResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &changeResponse, nil
}

func changeMyEmployerAccountPassword(url string, client *http.Client, token string, data *ChangeEmployerPasswordRequest) (int, *ChangeEmployerPasswordResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest("POST", url+"/api/v1/employers/accounts/me/change_password", strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var changeResponse ChangeEmployerPasswordResponse
	if err := json.NewDecoder(res.Body).Decode(&changeResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &changeResponse, nil
}

func TestUsersHandlersCreateEmployerAccount(t *testing.T) {
	ts, _, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()
	client := ts.Client()
	var adminToken string
	t.Run("prepare admin accounts and tokens", prepareAdminAccount(ts.URL, client, &adminToken))

	testCases := []struct {
		desc             string
		auth             string
		data             CreateEmployerAccountParams
		expectedResponse CreateEmployerAccountResponse
		status           int
	}{
		{
			desc: "Admin create user - success",
			data: CreateEmployerAccountParams{Email: "testuser@guillaume.com", Password: "Password1!"},
			auth: adminToken,
			expectedResponse: CreateEmployerAccountResponse{
				Result: CreateEmployerAccountResponseResult{
					Id: 1,
				},
			},
			status: http.StatusCreated,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := createEmployerAccount(ts.URL, client, tC.auth, "1", &tC.data)
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

func TestHandlersGetEmployerAccount(t *testing.T) {
	ts, _, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()
	client := ts.Client()
	var adminToken string
	var employerToken string
	t.Run("prepare admin accounts and tokens", prepareAdminAccount(ts.URL, client, &adminToken))
	t.Run("prepare user accounts and tokens", prepareEmployerAccount(ts.URL, client, &adminToken, &employerToken))

	testCases := []struct {
		desc             string
		id               string
		auth             string
		expectedResponse GetEmployerAccountResponse
		status           int
	}{
		{
			desc: "Admin get admin user - success",
			id:   "1",
			auth: adminToken,
			expectedResponse: GetEmployerAccountResponse{
				Result: GetEmployerAccountResponseResult{
					Id: 1, Email: "testadmin",
				},
			},
			status: http.StatusOK,
		},
		{
			desc:             "Admin get inexistent user - fail",
			id:               "300",
			auth:             adminToken,
			expectedResponse: GetEmployerAccountResponse{Error: "Admin Account not found"},
			status:           http.StatusNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := getEmployerAccount(ts.URL, client, tC.auth, tC.id)
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

func TestHandlersGetMyEmployerAccount(t *testing.T) {
	ts, _, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()
	client := ts.Client()
	var adminToken string
	var employerToken string
	t.Run("prepare admin accounts and tokens", prepareAdminAccount(ts.URL, client, &adminToken))
	t.Run("prepare user accounts and tokens", prepareEmployerAccount(ts.URL, client, &adminToken, &employerToken))

	testCases := []struct {
		desc             string
		auth             string
		expectedResponse GetEmployerAccountResponse
		status           int
	}{
		{
			desc: "Get own user - success",
			auth: employerToken,
			expectedResponse: GetEmployerAccountResponse{
				Result: GetEmployerAccountResponseResult{
					Id: 1, Email: "employee@testemployer.com",
				},
			},
			status: http.StatusOK,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := getMyEmployerAccount(ts.URL, client, tC.auth)
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

func TestUsersHandlersDeleteEmployerAccount(t *testing.T) {
	ts, _, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()
	client := ts.Client()
	var adminToken string
	var employerToken string
	t.Run("prepare admin accounts and tokens", prepareAdminAccount(ts.URL, client, &adminToken))
	t.Run("prepare user accounts and tokens", prepareEmployerAccount(ts.URL, client, &adminToken, &employerToken))

	testCases := []struct {
		desc             string
		method           string
		id               string
		auth             string
		expectedResponse DeleteEmployerAccountResponse
		status           int
	}{

		{
			desc:             "Admin delete non-existent user - failure",
			id:               "123",
			auth:             adminToken,
			expectedResponse: DeleteEmployerAccountResponse{Error: "Admin Account not found"},
			status:           http.StatusNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := deleteEmployerAccount(ts.URL, client, tC.auth, tC.id)
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

func TestHandlersChangeEmployerAccountPassword(t *testing.T) {
	ts, _, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()
	client := ts.Client()
	var adminToken string
	var employerToken string
	t.Run("prepare admin accounts and tokens", prepareAdminAccount(ts.URL, client, &adminToken))
	t.Run("prepare user accounts and tokens", prepareEmployerAccount(ts.URL, client, &adminToken, &employerToken))

	testCases := []struct {
		desc             string
		id               string
		data             ChangeEmployerPasswordRequest
		auth             string
		expectedResponse ChangeEmployerPasswordResponse
		status           int
	}{
		{
			desc:             "Admin change password with id - success",
			id:               "1",
			data:             ChangeEmployerPasswordRequest{Password: "newPassword1"},
			auth:             adminToken,
			expectedResponse: ChangeEmployerPasswordResponse{Result: EmployerPasswordResponseResult{Id: 1}},
			status:           http.StatusOK,
		},
		{
			desc:             "Admin change non-existing user password - failure",
			id:               "100",
			data:             ChangeEmployerPasswordRequest{Password: "newPassword1"},
			auth:             adminToken,
			expectedResponse: ChangeEmployerPasswordResponse{Error: "Admin Account not found"},
			status:           http.StatusNotFound,
		},
		{
			desc:             "Admin change password without password - failure",
			id:               "1",
			data:             ChangeEmployerPasswordRequest{},
			auth:             adminToken,
			expectedResponse: ChangeEmployerPasswordResponse{Error: "Password is required"},
			status:           http.StatusBadRequest,
		},
		{
			desc:             "Admin change password with bad password - failure",
			id:               "1",
			data:             ChangeEmployerPasswordRequest{Password: "password"},
			auth:             adminToken,
			expectedResponse: ChangeEmployerPasswordResponse{Error: "Password must have 8 or more characters, must include at least one capital letter, one lowercase letter, and either a number or a symbol."},
			status:           http.StatusBadRequest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := changeEmployerAccountPassword(ts.URL, client, tC.auth, tC.id, &tC.data)
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

func TestUsersHandlersChangeMyEmployerPassword(t *testing.T) {
	ts, _, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()
	client := ts.Client()
	var adminToken string
	var employerToken string
	t.Run("prepare admin accounts and tokens", prepareAdminAccount(ts.URL, client, &adminToken))
	t.Run("prepare user accounts and tokens", prepareEmployerAccount(ts.URL, client, &adminToken, &employerToken))

	testCases := []struct {
		desc             string
		data             ChangeEmployerPasswordRequest
		auth             string
		expectedResponse ChangeEmployerPasswordResponse
		status           int
	}{
		{
			desc:             "Admin change password with me - success",
			data:             ChangeEmployerPasswordRequest{Password: "newPassword1"},
			auth:             employerToken,
			expectedResponse: ChangeEmployerPasswordResponse{Result: EmployerPasswordResponseResult{Id: 1}},
			status:           http.StatusOK,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := changeMyEmployerAccountPassword(ts.URL, client, tC.auth, &tC.data)
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
