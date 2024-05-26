package handler

import (
	"net/http"
)

const (
	MAX_FILE_SIZE = 1024 * 1024 * 8 // 1MB
)

// POST /files
func (mux *ServerMux) handlePost(w http.ResponseWriter, r *http.Request) {
	// バリデーション
	if err := r.ParseMultipartForm(MAX_FILE_SIZE); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	var (
		err       error
		fileName  string
		title     string
		fileSize  int
		mimeType  string
		userEmail string
	)
	fileSrc, fileHeader, err := r.FormFile("upload-image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer fileSrc.Close()
	if len(fileHeader.Filename) == 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	fileName = fileHeader.Filename
	if fileHeader.Size == 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	fileSize = int(fileHeader.Size)
	if len(fileHeader.Header.Get("Content-Type")) == 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	mimeType = fileHeader.Header.Get("Content-Type")
	title = r.MultipartForm.Value["title"][0]
	if len(title) == 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	userEmail = r.MultipartForm.Value["email"][0]
	// email はもっとしっかりチェックしていい
	if len(userEmail) == 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// ユースケース
	if err := mux.fileUseCase.Post(
		r.Context(),
		fileSrc,
		fileName,
		title,
		fileSize,
		mimeType,
		userEmail,
	); err != nil {
		http.Error(w, "server error occurred", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
