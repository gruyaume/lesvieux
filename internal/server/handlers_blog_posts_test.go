package server_test

type UpdateJobPostParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

type UpdateJobPostResponseResult struct {
	ID int64 `json:"id"`
}

type UpdateJobPostResponse struct {
	Error  string                      `json:"error,omitempty"`
	Result UpdateJobPostResponseResult `json:"result"`
}

type ListJobPostsResponseResult []int

type ListJobPostsResponse struct {
	Error  string                     `json:"error,omitempty"`
	Result ListJobPostsResponseResult `json:"result"`
}

type GetJobPostResponseResult struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Status  string `json:"status,omitempty"`
}

type GetJobPostResponse struct {
	Error  string                   `json:"error,omitempty"`
	Result GetJobPostResponseResult `json:"result"`
}

type DeleteJobPostResponseResult struct {
	ID int64 `json:"id"`
}

type DeleteJobPostResponse struct {
	Error  string                      `json:"error,omitempty"`
	Result DeleteJobPostResponseResult `json:"result"`
}
