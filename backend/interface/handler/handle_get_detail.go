package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kawabatas/go-fileuploader/domain/model"
)

// GET /files/{id}
func (mux *ServerMux) handleGetDetail(w http.ResponseWriter, r *http.Request) {
	// バリデーション
	id := r.PathValue("id")
	if len(id) != 36 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// ユースケース
	file, err := mux.fileUseCase.GetDetail(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "server error occurred", http.StatusInternalServerError)
		return
	}

	// データ整形
	downloadURL, err := file.GetDownloadURL()
	if err != nil {
		http.Error(w, "server error occurred", http.StatusInternalServerError)
		return
	}
	res := &fileDetailResponse{
		Title:       file.Title,
		FileName:    file.FileName,
		Size:        file.Size,
		DownloadURL: downloadURL,
		CreatedAt:   file.CreatedAt,
		UserEmail:   file.User.Email,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "server error occurred", http.StatusInternalServerError)
		return
	}
}
