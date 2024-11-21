package server_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

type CreateEmployerParams struct {
	Name string `json:"name"`
}

type GetEmployerResponseResult struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type GetEmployerResponse struct {
	Result GetEmployerResponseResult `json:"result"`
	Error  string                    `json:"error,omitempty"`
}

type CreateEmployerResponseResult struct {
	Id int `json:"id"`
}

type CreateEmployerResponse struct {
	Result CreateEmployerResponseResult `json:"result"`
	Error  string                       `json:"error,omitempty"`
}

type DeleteEmployerResponseResult struct {
	Id int `json:"id"`
}

type DeleteEmployerResponse struct {
	Error  string                       `json:"error"`
	Result DeleteEmployerResponseResult `json:"result"`
}

func createEmployer(url string, client *http.Client, token string, data *CreateEmployerParams) (int, *CreateEmployerResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest("POST", url+"/api/v1/employers", strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var createResponse CreateEmployerResponse
	if err := json.NewDecoder(res.Body).Decode(&createResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &createResponse, nil
}

func getEmployer(url string, client *http.Client, token string, id string) (int, *GetEmployerResponse, error) {
	req, err := http.NewRequest("GET", url+"/api/v1/employers/"+id, nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var getResponse GetEmployerResponse
	if err := json.NewDecoder(res.Body).Decode(&getResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &getResponse, nil
}

func deleteEmployer(url string, client *http.Client, token string, id string) (int, *DeleteEmployerResponse, error) {
	req, err := http.NewRequest("DELETE", url+"/api/v1/employers/"+id, nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var deleteResponse DeleteEmployerResponse
	if err := json.NewDecoder(res.Body).Decode(&deleteResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &deleteResponse, nil
}

func TestHandlersCreateEmployers(t *testing.T) {
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
		data             CreateEmployerParams
		expectedResponse CreateEmployerResponse
		status           int
	}{
		{
			desc: "Create Employer - success",
			data: CreateEmployerParams{Name: "test employer"},
			auth: token,
			expectedResponse: CreateEmployerResponse{
				Result: CreateEmployerResponseResult{
					Id: 1,
				},
			},
			status: http.StatusCreated,
		},
		{
			desc: "Create Employer - failure (no name)",
			data: CreateEmployerParams{},
			auth: token,
			expectedResponse: CreateEmployerResponse{
				Error: "Name is required",
			},
			status: http.StatusBadRequest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := createEmployer(ts.URL, client, tC.auth, &tC.data)
			if err != nil {
				t.Fatalf("couldn't create Employer: %s", err)
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

func TestHandlersGetEmployers(t *testing.T) {
	ts, _, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()
	client := ts.Client()
	var token string
	t.Run("prepare user accounts and tokens", prepareAdminAccount(ts.URL, client, &token))
	t.Run("Create Employer", func(t *testing.T) {
		employer := &CreateEmployerParams{Name: "test employer"}
		statusCode, _, err := createEmployer(ts.URL, client, token, employer)
		if err != nil {
			t.Fatalf("couldn't create employer: %s", err)
		}
		if statusCode != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, statusCode)
		}
	})

	testCases := []struct {
		desc             string
		auth             string
		expectedResponse GetEmployerResponse
		status           int
	}{
		{
			desc: "Get Employer - success",
			auth: token,
			expectedResponse: GetEmployerResponse{
				Result: GetEmployerResponseResult{
					Id: 1, Name: "test employer",
				},
			},
			status: http.StatusOK,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := getEmployer(ts.URL, client, tC.auth, "1")
			if err != nil {
				t.Fatalf("couldn't get Employer: %s", err)
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

func TestHandlersDeleteEmployers(t *testing.T) {
	ts, _, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()
	client := ts.Client()
	var token string
	t.Run("prepare user accounts and tokens", prepareAdminAccount(ts.URL, client, &token))
	t.Run("Create Employer", func(t *testing.T) {
		employer := &CreateEmployerParams{Name: "test employer"}
		statusCode, _, err := createEmployer(ts.URL, client, token, employer)
		if err != nil {
			t.Fatalf("couldn't create employer: %s", err)
		}
		if statusCode != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, statusCode)
		}
	})

	testCases := []struct {
		desc             string
		auth             string
		id               string
		expectedResponse DeleteEmployerResponse
		status           int
	}{
		{
			desc: "Delete Employer - success",
			auth: token,
			id:   "1",
			expectedResponse: DeleteEmployerResponse{
				Result: DeleteEmployerResponseResult{
					Id: 1,
				},
			},
			status: http.StatusAccepted,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statusCode, resp, err := deleteEmployer(ts.URL, client, tC.auth, tC.id)
			if err != nil {
				t.Fatalf("couldn't delete Employer: %s", err)
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
