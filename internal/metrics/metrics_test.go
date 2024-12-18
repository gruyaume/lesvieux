package metrics_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruyaume/lesvieux/internal/db"
	"github.com/gruyaume/lesvieux/internal/metrics"
)

// TestPrometheusHandler tests that the Prometheus metrics handler responds correctly to an HTTP request.
func TestPrometheusHandler(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "db.sqlite3")

	dbQueries, err := db.Initialize(dbPath)
	if err != nil {
		t.Fatal(err)
	}

	m := metrics.NewMetricsSubsystem(dbQueries)

	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	recorder := httptest.NewRecorder()
	m.Handler.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	if recorder.Body.String() == "" {
		t.Errorf("handler returned an empty body")
	}
	if !strings.Contains(recorder.Body.String(), "go_goroutines") {
		t.Errorf("expected 'go_goroutines' in the metrics output, but it was missing")
	}
}

// TestMetrics tests some of the metrics that we currently collect.
func TestMetrics(t *testing.T) {
	dbQueries, err := db.Initialize(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	createJobPost1Params := db.CreateJobPostParams{
		Title:     "my title 1",
		Content:   "my content 1",
		CreatedAt: "creation time 1",
	}
	createJobPost2Params := db.CreateJobPostParams{
		Title:     "my title 2",
		Content:   "my content 2",
		CreatedAt: "creation time 2",
	}
	_, err = dbQueries.CreateJobPost(context.Background(), createJobPost1Params)
	if err != nil {
		t.Fatalf("couldn't create test job post: %s", err)
	}
	_, err = dbQueries.CreateJobPost(context.Background(), createJobPost2Params)
	if err != nil {
		t.Fatalf("couldn't create test job post: %s", err)
	}

	m := metrics.NewMetricsSubsystem(dbQueries)
	jobPosts, err := dbQueries.ListJobPosts(context.Background())
	if err != nil {
		t.Fatalf("couldn't list job posts: %s", err)
	}
	m.GenerateMetrics(jobPosts)

	request, _ := http.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()
	m.Handler.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	if recorder.Body.String() == "" {
		t.Errorf("handler returned an empty body")
	}
	for _, line := range strings.Split(recorder.Body.String(), "\n") {
		if strings.Contains(line, "job_posts_total ") && !strings.HasPrefix(line, "#") {
			if !strings.HasSuffix(line, "2") {
				t.Errorf("job_posts_total expected to receive 2")
			}
		}
	}
}
