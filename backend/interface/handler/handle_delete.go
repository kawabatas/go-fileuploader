package handler

import (
	"errors"
	"net/http"

	"github.com/kawabatas/go-fileuploader/domain/model"
)

// DELETE /files/{id}
func (mux *ServerMux) handleDelete(w http.ResponseWriter, r *http.Request) {
	// バリデーション
	id := r.PathValue("id")
	if len(id) != 36 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// ユースケース
	if err := mux.fileUseCase.Delete(r.Context(), id); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "server error occurred", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
