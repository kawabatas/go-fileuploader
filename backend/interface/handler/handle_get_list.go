package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const (
	FILE_COUNT_PER_PAGE = 10
)

// GET /
func (mux *ServerMux) handleGetList(w http.ResponseWriter, r *http.Request) {
	offset := 0
	offsetStr := r.URL.Query().Get("offset")
	if len(offsetStr) > 0 {
		num, err := strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		offset = num
	}
	limit := offset + FILE_COUNT_PER_PAGE

	// ユースケース limit+1 取得して、ページャー判定
	files, err := mux.fileUseCase.GetList(r.Context(), offset, limit+1)
	if err != nil {
		http.Error(w, "server error occurred", http.StatusInternalServerError)
		return
	}

	// データ整形
	resFiles := make([]*fileResponse, len(files))
	for i, f := range files {
		resFiles[i] = &fileResponse{
			ID:    f.ID,
			Title: f.Title,
		}
	}
	res := &fileListResponse{
		Files: resFiles,
		Prev:  offset != 0,
		Next:  len(files) > FILE_COUNT_PER_PAGE,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "server error occurred", http.StatusInternalServerError)
		return
	}
}
