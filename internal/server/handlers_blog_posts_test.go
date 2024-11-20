package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

type createMyBlogPostResponseResult struct {
	ID int64 `json:"id"`
}

type createMyBlogPostResponse struct {
	Error  string                         `json:"error,omitempty"`
	Result createMyBlogPostResponseResult `json:"result"`
}

type UpdateBlogPostParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

type UpdateBlogPostResponseResult struct {
	ID int64 `json:"id"`
}

type UpdateBlogPostResponse struct {
	Error  string                       `json:"error,omitempty"`
	Result UpdateBlogPostResponseResult `json:"result"`
}

type ListBlogPostsResponseResult []int

type ListBlogPostsResponse struct {
	Error  string                      `json:"error,omitempty"`
	Result ListBlogPostsResponseResult `json:"result"`
}

type GetBlogPostResponseResult struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Status  string `json:"status,omitempty"`
}

type GetBlogPostResponse struct {
	Error  string                    `json:"error,omitempty"`
	Result GetBlogPostResponseResult `json:"result"`
}

type DeleteBlogPostResponseResult struct {
	ID int64 `json:"id"`
}

type DeleteBlogPostResponse struct {
	Error  string                       `json:"error,omitempty"`
	Result DeleteBlogPostResponseResult `json:"result"`
}

func listBlogPosts(url string, client *http.Client, token string) (int, *ListBlogPostsResponse, error) {
	req, err := http.NewRequest("GET", url+"/api/v1/posts", nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var listResponse ListBlogPostsResponse
	if err := json.NewDecoder(res.Body).Decode(&listResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &listResponse, nil
}

func getBlogPost(url string, client *http.Client, token string, id int) (int, *GetBlogPostResponse, error) {
	req, err := http.NewRequest("GET", url+"/api/v1/posts/"+strconv.Itoa(id), nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var getResponse GetBlogPostResponse
	if err := json.NewDecoder(res.Body).Decode(&getResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &getResponse, nil
}

func listMyBlogPosts(url string, client *http.Client, token string) (int, *ListBlogPostsResponse, error) {
	req, err := http.NewRequest("GET", url+"/api/v1/me/posts", nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var listResponse ListBlogPostsResponse
	if err := json.NewDecoder(res.Body).Decode(&listResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &listResponse, nil
}

func listPublicBlogPosts(url string, client *http.Client) (int, *ListBlogPostsResponse, error) {
	req, err := http.NewRequest("GET", url+"/api/v1/published_posts", nil)
	if err != nil {
		return 0, nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var listResponse ListBlogPostsResponse
	if err := json.NewDecoder(res.Body).Decode(&listResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &listResponse, nil
}

func getPublicBlogPost(url string, client *http.Client, token string, id int) (int, *GetBlogPostResponse, error) {
	req, err := http.NewRequest("GET", url+"/api/v1/published_posts/"+strconv.Itoa(id), nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var getResponse GetBlogPostResponse
	if err := json.NewDecoder(res.Body).Decode(&getResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &getResponse, nil
}

func getMyBlogPost(url string, client *http.Client, token string, id int) (int, *GetBlogPostResponse, error) {
	req, err := http.NewRequest("GET", url+"/api/v1/me/posts/"+strconv.Itoa(id), nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var getResponse GetBlogPostResponse
	if err := json.NewDecoder(res.Body).Decode(&getResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &getResponse, nil
}

func createMyBlogPost(url string, client *http.Client, token string) (int, *createMyBlogPostResponse, error) {
	req, err := http.NewRequest("POST", url+"/api/v1/me/posts", nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var createResponse createMyBlogPostResponse
	if err := json.NewDecoder(res.Body).Decode(&createResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &createResponse, nil
}

func updateMyBlogPost(url string, client *http.Client, token string, id int, blogPost UpdateBlogPostParams) (int, *UpdateBlogPostResponse, error) {
	body, err := json.Marshal(blogPost)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest("PUT", url+"/api/v1/me/posts/"+strconv.Itoa(id), bytes.NewReader(body))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var updateResponse UpdateBlogPostResponse
	if err := json.NewDecoder(res.Body).Decode(&updateResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &updateResponse, nil
}

func deleteMyBlogPost(url string, client *http.Client, token string, id int) (int, error) {
	req, err := http.NewRequest("DELETE", url+"/api/v1/me/posts/"+strconv.Itoa(id), nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	return res.StatusCode, nil
}

func TestPublicPostsHandlers(t *testing.T) {
	ts, _, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()
	client := ts.Client()

	var adminToken string
	var nonAdminToken string
	t.Run("prepare accounts and tokens", prepareUserAccounts(ts.URL, client, &adminToken, &nonAdminToken))

	t.Run("List public blog posts - 0", func(t *testing.T) {
		statusCode, response, err := listPublicBlogPosts(ts.URL, client)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}

		if response.Error != "" {
			t.Fatalf("expected error %q, got %q", "", response.Error)
		}

		if len(response.Result) != 0 {
			t.Fatalf("expected result %v, got %v", []int{}, response.Result)
		}
	})

	t.Run("Create draft blog post", func(t *testing.T) {
		statusCode, response, err := createMyBlogPost(ts.URL, client, adminToken)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, statusCode)
		}

		if response.Result.ID != 1 {
			t.Fatalf("expected id %d, got %d", 1, response.Result.ID)
		}
	})

	t.Run("List public blog posts - 0", func(t *testing.T) {
		statusCode, response, err := listPublicBlogPosts(ts.URL, client)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}

		if response.Error != "" {
			t.Fatalf("expected error %q, got %q", "", response.Error)
		}

		if len(response.Result) != 0 {
			t.Fatalf("expected result %v, got %v", []int{}, response.Result)
		}
	})

	t.Run("Update blog post to published", func(t *testing.T) {
		blogPost := UpdateBlogPostParams{
			Title:   "Test Title",
			Content: "Test Content",
			Status:  "published",
		}
		statusCode, response, err := updateMyBlogPost(ts.URL, client, adminToken, 1, blogPost)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}

		if response.Result.ID != 1 {
			t.Fatalf("expected id %d, got %d", 1, response.Result.ID)
		}
	})

	t.Run("List public blog posts - 1", func(t *testing.T) {
		statusCode, response, err := listPublicBlogPosts(ts.URL, client)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}

		if response.Error != "" {
			t.Fatalf("expected error %q, got %q", "", response.Error)
		}

		if len(response.Result) != 1 {
			t.Fatalf("expected result %v, got %v", []int{1}, response.Result)
		}
	})

	t.Run("Get blog post", func(t *testing.T) {
		expectedBlogPost := GetBlogPostResponseResult{
			Title:   "Test Title",
			Content: "Test Content",
			Status:  "published",
		}
		statusCode, response, err := getPublicBlogPost(ts.URL, client, adminToken, 1)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}

		if response.Error != "" {
			t.Fatalf("expected error %q, got %q", "", response.Error)
		}

		if response.Result.Title != expectedBlogPost.Title {
			t.Fatalf("expected title %q, got %q", expectedBlogPost.Title, response.Result.Title)
		}

		if response.Result.Content != expectedBlogPost.Content {
			t.Fatalf("expected content %q, got %q", expectedBlogPost.Content, response.Result.Content)
		}

		if response.Result.Status != expectedBlogPost.Status {
			t.Fatalf("expected status %q, got %q", expectedBlogPost.Status, response.Result.Status)
		}
	})
}

func TestMyBlogPostsHandlers(t *testing.T) {
	ts, _, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()
	client := ts.Client()

	var adminToken string
	var nonAdminToken string
	t.Run("prepare accounts and tokens", prepareUserAccounts(ts.URL, client, &adminToken, &nonAdminToken))

	t.Run("Author List author blog posts", func(t *testing.T) {
		statusCode, response, err := listMyBlogPosts(ts.URL, client, nonAdminToken)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}

		if response.Error != "" {
			t.Fatalf("expected error %q, got %q", "", response.Error)
		}

		if len(response.Result) != 0 {
			t.Fatalf("expected result %v, got %v", []int{}, response.Result)
		}
	})

	t.Run("Create blog post", func(t *testing.T) {
		statusCode, _, err := createMyBlogPost(ts.URL, client, nonAdminToken)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, statusCode)
		}
	})

	t.Run("List author blog posts", func(t *testing.T) {
		statusCode, response, err := listMyBlogPosts(ts.URL, client, nonAdminToken)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}

		if response.Error != "" {
			t.Fatalf("expected error %q, got %q", "", response.Error)
		}

		if len(response.Result) != 1 {
			t.Fatalf("expected result %v, got %v", []int{1}, response.Result)
		}
	})

	t.Run("Get blog post", func(t *testing.T) {
		statusCode, response, err := getMyBlogPost(ts.URL, client, nonAdminToken, 1)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}

		if response.Error != "" {
			t.Fatalf("expected error %q, got %q", "", response.Error)
		}

		if response.Result.Status != "draft" {
			t.Fatalf("expected status %q, got %q", "draft", response.Result.Status)
		}
	})

	t.Run("Update blog post", func(t *testing.T) {
		blogPost := UpdateBlogPostParams{
			Title:   "Test Title",
			Content: "Test Content",
			Status:  "published",
		}
		statusCode, response, err := updateMyBlogPost(ts.URL, client, nonAdminToken, 1, blogPost)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}

		if response.Result.ID != 1 {
			t.Fatalf("expected id %d, got %d", 1, response.Result.ID)
		}
	})

	t.Run("Get blog post", func(t *testing.T) {
		statusCode, response, err := getMyBlogPost(ts.URL, client, nonAdminToken, 1)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}

		if response.Error != "" {
			t.Fatalf("expected error %q, got %q", "", response.Error)
		}

		if response.Result.Title != "Test Title" {
			t.Fatalf("expected title %q, got %q", "Test Title", response.Result.Title)
		}

		if response.Result.Content != "Test Content" {
			t.Fatalf("expected content %q, got %q", "Test Content", response.Result.Content)
		}

		if response.Result.Status != "published" {
			t.Fatalf("expected status %q, got %q", "published", response.Result.Status)
		}
	})

	t.Run("Create another blog post", func(t *testing.T) {
		statusCode, _, err := createMyBlogPost(ts.URL, client, nonAdminToken)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, statusCode)
		}
	})

	t.Run("List author blog posts", func(t *testing.T) {
		statusCode, response, err := listMyBlogPosts(ts.URL, client, nonAdminToken)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}

		if response.Error != "" {
			t.Fatalf("expected error %q, got %q", "", response.Error)
		}

		if len(response.Result) != 2 {
			t.Fatalf("expected result %v, got %v", []int{1, 2}, response.Result)
		}
	})

	t.Run("Delete first blog post", func(t *testing.T) {
		statusCode, err := deleteMyBlogPost(ts.URL, client, nonAdminToken, 1)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusAccepted {
			t.Fatalf("expected status %d, got %d", http.StatusAccepted, statusCode)
		}
	})

	t.Run("Get deleted blog post", func(t *testing.T) {
		statusCode, _, err := getMyBlogPost(ts.URL, client, adminToken, 1)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, statusCode)
		}
	})

	t.Run("List blog posts", func(t *testing.T) {
		statusCode, response, err := listMyBlogPosts(ts.URL, client, nonAdminToken)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}

		if response.Error != "" {
			t.Fatalf("expected error %q, got %q", "", response.Error)
		}

		if len(response.Result) != 1 {
			t.Fatalf("expected result %v, got %v", []int{2}, response.Result)
		}
	})
}

func TestAdminBlogPostsHandlers(t *testing.T) {
	ts, _, err := setupServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Close()
	client := ts.Client()

	var adminToken string
	var nonAdminToken string
	t.Run("prepare accounts and tokens", prepareUserAccounts(ts.URL, client, &adminToken, &nonAdminToken))

	t.Run("List blog posts", func(t *testing.T) {
		statusCode, response, err := listBlogPosts(ts.URL, client, adminToken)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}

		if response.Error != "" {
			t.Fatalf("expected error %q, got %q", "", response.Error)
		}

		if len(response.Result) != 0 {
			t.Fatalf("expected result %v, got %v", []int{}, response.Result)
		}
	})

	t.Run("Create blog post - nonAdmin", func(t *testing.T) {
		statusCode, resp, err := createMyBlogPost(ts.URL, client, nonAdminToken)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, statusCode)
		}

		if resp.Result.ID != 1 {
			t.Fatalf("expected id %d, got %d", 1, resp.Result.ID)
		}
	})

	t.Run("Create blog post - admin", func(t *testing.T) {
		statusCode, resp, err := createMyBlogPost(ts.URL, client, adminToken)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, statusCode)
		}

		if resp.Result.ID != 2 {
			t.Fatalf("expected id %d, got %d", 2, resp.Result.ID)
		}
	})

	t.Run("Update blog post - admin", func(t *testing.T) {
		blogPost := UpdateBlogPostParams{
			Title:   "Test Title",
			Content: "Test Content",
			Status:  "published",
		}
		statusCode, response, err := updateMyBlogPost(ts.URL, client, adminToken, 2, blogPost)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}

		if response.Result.ID != 2 {
			t.Fatalf("expected id %d, got %d", 2, response.Result.ID)
		}
	})

	t.Run("List blog posts", func(t *testing.T) {
		statusCode, response, err := listBlogPosts(ts.URL, client, adminToken)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}

		if response.Error != "" {
			t.Fatalf("expected error %q, got %q", "", response.Error)
		}

		if len(response.Result) != 2 {
			t.Fatalf("expected result %v, got %v", []int{1, 2}, response.Result)
		}
	})

	t.Run("Get admin blog post", func(t *testing.T) {
		statusCode, response, err := getBlogPost(ts.URL, client, adminToken, 2)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}

		if response.Error != "" {
			t.Fatalf("expected error %q, got %q", "", response.Error)
		}

		if response.Result.Title != "Test Title" {
			t.Fatalf("expected title %q, got %q", "Test Title", response.Result.Title)
		}

		if response.Result.Content != "Test Content" {
			t.Fatalf("expected content %q, got %q", "Test Content", response.Result.Content)
		}

		if response.Result.Status != "published" {
			t.Fatalf("expected status %q, got %q", "published", response.Result.Status)
		}
	})

	t.Run("Get non-admin blog post", func(t *testing.T) {
		statusCode, response, err := getBlogPost(ts.URL, client, adminToken, 1)
		if err != nil {
			t.Fatal(err)
		}

		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}

		if response.Error != "" {
			t.Fatalf("expected error %q, got %q", "", response.Error)
		}

		if response.Result.Status != "draft" {
			t.Fatalf("expected status %q, got %q", "draft", response.Result.Status)
		}
	})
}
