package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/IZO-Ong/gridgo/internal/db"
	"github.com/IZO-Ong/gridgo/internal/models"
)

func HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "USERNAME_REQUIRED", http.StatusBadRequest)
		return
	}

	var user models.User
	err := db.DB.Preload("Mazes").Where("username = ?", username).First(&user).Error
	
	if err != nil {
		http.Error(w, "USER_NOT_FOUND", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"username":   user.Username,
		"created_at": user.CreatedAt,
		"mazes":      user.Mazes,
		"stats": map[string]interface{}{
			"total_mazes": len(user.Mazes),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}