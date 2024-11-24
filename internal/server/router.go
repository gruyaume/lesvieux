package server

import (
	"net/http"

	"github.com/gruyaume/lesvieux/internal/metrics"
)

func NewLesVieuxRouter(config *HandlerConfig) http.Handler {
	apiV1Router := http.NewServeMux()

	// No Auth
	apiV1Router.HandleFunc("POST /employers/login", EmployersLogin(config))
	apiV1Router.HandleFunc("POST /admin/login", AdminLogin(config))
	apiV1Router.HandleFunc("GET /status", GetStatus(config))
	apiV1Router.HandleFunc("GET /posts", ListJobPosts(config))

	// Admin or First User
	apiV1Router.HandleFunc("POST /admin/accounts", adminOrFirstUser(config.JWTSecret, config.DBQueries, CreateAdminAccount(config)))

	// Admin Only
	apiV1Router.HandleFunc("GET /posts/{post_id}", adminOnly(config.JWTSecret, GetJobPost(config)))
	apiV1Router.HandleFunc("POST /employers", adminOnly(config.JWTSecret, CreateEmployer(config)))
	apiV1Router.HandleFunc("GET /employers", adminOnly(config.JWTSecret, ListEmployers(config)))
	apiV1Router.HandleFunc("GET /employers/{employer_id}", adminOnly(config.JWTSecret, GetEmployer(config)))
	apiV1Router.HandleFunc("DELETE /employers/{employer_id}", adminOnly(config.JWTSecret, DeleteEmployer(config)))
	apiV1Router.HandleFunc("GET /employers/{employer_id}/accounts", adminOnly(config.JWTSecret, ListEmployerAccounts(config)))
	apiV1Router.HandleFunc("POST /employers/{employer_id}/accounts", adminOnly(config.JWTSecret, CreateEmployerAccount(config)))
	apiV1Router.HandleFunc("GET /employers/{employer_id}/accounts/{account_id}", adminOnly(config.JWTSecret, GetEmployerAccount(config)))
	apiV1Router.HandleFunc("DELETE /employers/{employer_id}/accounts/{account_id}", adminOnly(config.JWTSecret, DeleteEmployerAccount(config)))
	apiV1Router.HandleFunc("POST /employers/{employer_id}/accounts/{account_id}/change_password", adminOnly(config.JWTSecret, ChangeEmployerAccountPassword(config)))
	apiV1Router.HandleFunc("GET /admin/accounts", adminOnly(config.JWTSecret, ListAdminAccounts(config)))
	apiV1Router.HandleFunc("GET /admin/accounts/{account_id}", adminOnly(config.JWTSecret, GetAdminAccount(config)))
	apiV1Router.HandleFunc("DELETE /admin/accounts/{account_id}", adminOnly(config.JWTSecret, DeleteAdminAccount(config)))
	apiV1Router.HandleFunc("POST /admin/accounts/{account_id}/change_password", adminOnly(config.JWTSecret, ChangeAdminAccountPassword(config)))

	// Admin (Me) Only
	apiV1Router.HandleFunc("GET /employers/accounts/me", Me(config.JWTSecret, GetMyEmployerAccount(config)))
	apiV1Router.HandleFunc("POST /employers/accounts/me/change_password", Me(config.JWTSecret, ChangeMyEmployerAccountPassword(config)))
	apiV1Router.HandleFunc("GET /admin/accounts/me", Me(config.JWTSecret, GetMyAdminAccount(config)))
	apiV1Router.HandleFunc("POST /admin/accounts/me/change_password", Me(config.JWTSecret, ChangeMyAdminAccountPassword(config)))

	frontendHandler := newFrontendFileServer()

	router := http.NewServeMux()
	router.HandleFunc("GET /status", GetStatus(config))
	m := metrics.NewMetricsSubsystem(config.DBQueries)
	router.Handle("/metrics", m.Handler)
	ctx := loggingMiddlewareContext{}
	apiMiddlewareStack := createMiddlewareStack(
		metricsMiddleware(m),
		loggingMiddleware(&ctx),
	)
	metricsMiddlewareStack := createMiddlewareStack(
		metricsMiddleware(m),
	)
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMiddlewareStack(apiV1Router)))
	router.Handle("/", metricsMiddlewareStack(frontendHandler))

	return router
}
