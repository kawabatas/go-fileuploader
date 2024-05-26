package handler

import "time"

type fileResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type fileListResponse struct {
	Files []*fileResponse `json:"files"`
	Prev  bool            `json:"prev"`
	Next  bool            `json:"next"`
}

type fileDetailResponse struct {
	Title       string    `json:"title"`
	FileName    string    `json:"filename"`
	Size        int       `json:"size"`
	DownloadURL string    `json:"download_url"`
	CreatedAt   time.Time `json:"created_at"`
}
