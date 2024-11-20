package server

import (
	"net/http"

	"github.com/gruyaume/lesvieux/internal/metrics"
)

func NewHebdoRouter(config *HandlerConfig) http.Handler {
	apiV1Router := http.NewServeMux()

	// No Auth
	apiV1Router.HandleFunc("GET /published_posts", ListPublicBlogPosts(config))
	apiV1Router.HandleFunc("GET /published_posts/{id}", GetPublicBlogPost(config))
	apiV1Router.HandleFunc("POST /login", Login(config))
	apiV1Router.HandleFunc("GET /status", GetStatus(config))

	// Admin or First User
	apiV1Router.HandleFunc("POST /accounts", adminOrFirstUser(config.JWTSecret, config.DBQueries, CreateAccount(config)))

	// Admin Only
	apiV1Router.HandleFunc("GET /posts", adminOnly(config.JWTSecret, ListBlogPosts(config)))
	apiV1Router.HandleFunc("GET /posts/{post_id}", adminOnly(config.JWTSecret, GetBlogPost(config)))
	apiV1Router.HandleFunc("DELETE /posts/{post_id}", adminOnly(config.JWTSecret, DeleteBlogPost(config)))
	apiV1Router.HandleFunc("GET /accounts", adminOnly(config.JWTSecret, ListAccounts(config)))
	apiV1Router.HandleFunc("GET /accounts/{user_id}", adminOnly(config.JWTSecret, GetAccount(config)))
	apiV1Router.HandleFunc("POST /accounts/{user_id}/change_password", adminOnly(config.JWTSecret, ChangeAccountPassword(config)))
	apiV1Router.HandleFunc("DELETE /accounts/{user_id}", adminOnly(config.JWTSecret, DeleteAccount(config)))

	// Author (me) Only
	apiV1Router.HandleFunc("GET /me", Me(config.JWTSecret, GetMyAccount(config)))
	apiV1Router.HandleFunc("POST /me/change_password", Me(config.JWTSecret, ChangeMyAccountPassword(config)))
	apiV1Router.HandleFunc("GET /me/posts", Me(config.JWTSecret, ListMyBlogPosts(config)))
	apiV1Router.HandleFunc("POST /me/posts", Me(config.JWTSecret, CreateMyBlogPost(config)))
	apiV1Router.HandleFunc("GET /me/posts/{post_id}", Me(config.JWTSecret, GetMyBlogPost(config)))
	apiV1Router.HandleFunc("PUT /me/posts/{post_id}", Me(config.JWTSecret, UpdateMyBlogPost(config)))
	apiV1Router.HandleFunc("DELETE /me/posts/{post_id}", Me(config.JWTSecret, DeleteMyBlogPost(config)))

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
