package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

type CreateJobPostResponse struct {
	ID int64 `json:"id"`
}

type UpdateJobPostParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

type UpdateJobPostResponse struct {
	ID int64 `json:"id"`
}

type GetJobPostResponse struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	Author    string `json:"author"`
}

func ListJobPosts(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobPosts, err := env.DBQueries.ListJobPosts(context.Background())
		if err != nil {
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		ids := make([]int64, 0, len(jobPosts))
		for _, post := range jobPosts {
			ids = append(ids, post.ID)
		}

		err = writeJSON(w, ids)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

func GetJobPost(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("post_id")
		idInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "id must be an integer")
			return
		}

		jobPost, err := env.DBQueries.GetJobPost(context.Background(), idInt64)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "Job Post not found")
				return
			}
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		// Get the author of the job post
		author, err := env.DBQueries.GetAccount(context.Background(), jobPost.AccountID)
		if err != nil {
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		jobPostResponse := GetJobPostResponse{
			ID:        jobPost.ID,
			Title:     jobPost.Title,
			Content:   jobPost.Content,
			Status:    jobPost.Status,
			CreatedAt: jobPost.CreatedAt,
			Author:    author.Username,
		}

		w.WriteHeader(http.StatusOK)
		err = writeJSON(w, jobPostResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

func DeleteJobPost(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("post_id")
		idInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "id must be an integer")
			return
		}
		err = env.DBQueries.DeleteJobPost(context.Background(), idInt64)
		if err != nil {
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}
