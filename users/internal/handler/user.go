package handler

import (
	"encoding/json"
	"net/http"
)

func (h *UserHandler) getAllUsers(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    users, err := h.userService.GetUsers(ctx, r.URL.Query())
    if err != nil {
        http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(users); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
}