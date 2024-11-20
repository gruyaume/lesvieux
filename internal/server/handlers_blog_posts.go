package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gruyaume/lesvieux/internal/db"
)

type CreateBlogPostResponse struct {
	ID int64 `json:"id"`
}

type UpdateBlogPostParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

type UpdateBlogPostResponse struct {
	ID int64 `json:"id"`
}

type GetBlogPostResponse struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	Author    string `json:"author"`
}

func generateDate() string {
	now := time.Now()
	return now.Format(time.RFC3339)
}

func ListPublicBlogPosts(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		blogPosts, err := env.DBQueries.ListPublicBlogPosts(context.Background())
		if err != nil {
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		ids := make([]int64, 0, len(blogPosts))
		for _, post := range blogPosts {
			ids = append(ids, post.ID)
		}

		err = writeJSON(w, ids)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

func GetPublicBlogPost(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		idInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "id must be an integer")
			return
		}

		blogPost, err := env.DBQueries.GetPublicBlogPost(context.Background(), idInt64)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "Blog Post not found")
				return
			}
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		// Get the author of the blog post
		author, err := env.DBQueries.GetAccount(context.Background(), blogPost.AccountID)
		if err != nil {
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		blogPostResponse := GetBlogPostResponse{
			ID:        blogPost.ID,
			Title:     blogPost.Title,
			Content:   blogPost.Content,
			Status:    blogPost.Status,
			CreatedAt: blogPost.CreatedAt,
			Author:    author.Username,
		}

		w.WriteHeader(http.StatusOK)
		err = writeJSON(w, blogPostResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

func ListBlogPosts(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		blogPosts, err := env.DBQueries.ListBlogPosts(context.Background())
		if err != nil {
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		ids := make([]int64, 0, len(blogPosts))
		for _, post := range blogPosts {
			ids = append(ids, post.ID)
		}

		err = writeJSON(w, ids)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

func ListMyBlogPosts(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(userIDKey).(int64)
		if !ok {
			writeError(w, http.StatusUnauthorized, "Unauthorized - user id not found")
			return
		}

		blogPosts, err := env.DBQueries.ListBlogPostsByAccount(context.Background(), userID)
		if err != nil {
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		ids := make([]int64, 0, len(blogPosts))
		for _, post := range blogPosts {
			ids = append(ids, post.ID)
		}

		err = writeJSON(w, ids)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

func GetMyBlogPost(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("post_id")
		idInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "id must be an integer")
			return
		}

		userID, ok := r.Context().Value(userIDKey).(int64)
		if !ok {
			writeError(w, http.StatusUnauthorized, "Unauthorized - user id not found")
			return
		}

		blogPost, err := env.DBQueries.GetBlogPost(context.Background(), idInt64)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "Blog Post not found")
				return
			}
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		if blogPost.AccountID != userID {
			writeError(w, http.StatusForbidden, "forbidden: admin or user access required")
			return
		}

		// Get the author of the blog post
		author, err := env.DBQueries.GetAccount(context.Background(), blogPost.AccountID)
		if err != nil {
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		blogPostResponse := GetBlogPostResponse{
			ID:        blogPost.ID,
			Title:     blogPost.Title,
			Content:   blogPost.Content,
			Status:    blogPost.Status,
			CreatedAt: blogPost.CreatedAt,
			Author:    author.Username,
		}

		w.WriteHeader(http.StatusOK)
		err = writeJSON(w, blogPostResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

// CreateMyBlogPost creates a new Blog Post , and returns the id of the created row
func CreateMyBlogPost(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(userIDKey).(int64)
		if !ok {
			writeError(w, http.StatusUnauthorized, "Unauthorized - user id not found")
			return
		}
		createdAt := generateDate()
		blogPost := db.CreateBlogPostParams{
			CreatedAt: createdAt,
			Status:    "draft",
			AccountID: int64(userID),
		}
		dbBlogPost, err := env.DBQueries.CreateBlogPost(context.Background(), blogPost)
		if err != nil {
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		w.WriteHeader(http.StatusCreated)
		response := CreateBlogPostResponse{ID: dbBlogPost.ID}
		err = writeJSON(w, response)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

func UpdateMyBlogPost(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("post_id")
		idInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "id must be an integer")
			return
		}

		userID, ok := r.Context().Value(userIDKey).(int64)
		if !ok {
			writeError(w, http.StatusUnauthorized, "Unauthorized - user id not found")
			return
		}

		var updateBlogPostParams UpdateBlogPostParams
		if err := json.NewDecoder(r.Body).Decode(&updateBlogPostParams); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}

		blogPost, err := env.DBQueries.GetBlogPost(context.Background(), idInt64)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "Blog Post not found")
				return
			}
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		if blogPost.AccountID != userID {
			writeError(w, http.StatusForbidden, "forbidden: admin or user access required")
			return
		}
		blogPostUpdate := db.UpdateBlogPostParams{
			ID:      idInt64,
			Title:   updateBlogPostParams.Title,
			Content: updateBlogPostParams.Content,
			Status:  updateBlogPostParams.Status,
		}

		err = env.DBQueries.UpdateBlogPost(context.Background(), blogPostUpdate)
		if err != nil {
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		w.WriteHeader(http.StatusOK)
		response := UpdateBlogPostResponse{ID: idInt64}
		err = writeJSON(w, response)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

func GetBlogPost(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("post_id")
		idInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "id must be an integer")
			return
		}

		blogPost, err := env.DBQueries.GetBlogPost(context.Background(), idInt64)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "Blog Post not found")
				return
			}
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		// Get the author of the blog post
		author, err := env.DBQueries.GetAccount(context.Background(), blogPost.AccountID)
		if err != nil {
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		blogPostResponse := GetBlogPostResponse{
			ID:        blogPost.ID,
			Title:     blogPost.Title,
			Content:   blogPost.Content,
			Status:    blogPost.Status,
			CreatedAt: blogPost.CreatedAt,
			Author:    author.Username,
		}

		w.WriteHeader(http.StatusOK)
		err = writeJSON(w, blogPostResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

func DeleteMyBlogPost(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("post_id")
		idInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "id must be an integer")
			return
		}

		userID, ok := r.Context().Value(userIDKey).(int64)
		if !ok {
			writeError(w, http.StatusUnauthorized, "Unauthorized - user id not found")
			return
		}

		blogPost, err := env.DBQueries.GetBlogPost(context.Background(), idInt64)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "Blog Post not found")
				return
			}
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		if blogPost.AccountID != userID {
			writeError(w, http.StatusForbidden, "forbidden: admin or user access required")
			return
		}

		err = env.DBQueries.DeleteBlogPost(context.Background(), idInt64)
		if err != nil {
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

func DeleteBlogPost(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("post_id")
		idInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "id must be an integer")
			return
		}
		err = env.DBQueries.DeleteBlogPost(context.Background(), idInt64)
		if err != nil {
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}
