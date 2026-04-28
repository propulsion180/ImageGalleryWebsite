package handlers

import (
	"encoding/json"
	"gallery-server/db"
	"gallery-server/models"
	//	"io"
	"log/slog"
	"net/http"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var signup_data models.LoginData // Reused for signup
	err := json.NewDecoder(r.Body).Decode(&signup_data)
	if err != nil {
		slog.Error("Failed to decode JSON data for signup", "status", http.StatusBadRequest, "Error", err.Error())
		http.Error(w, "Invalid Request payload", http.StatusBadRequest)
		return
	}

	conn := db.ConnectDB()
	defer conn.Close()

	worked, err := db.AddUser(conn, &models.User{Username: signup_data.Username, Password: signup_data.Password, Admin: false})

	if !worked {
		slog.Error("Failed to add user", "status", http.StatusUnauthorized)
		http.Error(w, "Failed to add user", http.StatusUnauthorized)
		return
	}

	if err != nil {
		slog.Error("Error occured while adding user", "status", http.StatusInternalServerError)
		http.Error(w, "Error caught while add user", http.StatusInternalServerError)
		return
	}

	slog.Info("User added sucessfully", "status", http.StatusOK)
	w.WriteHeader(http.StatusOK)
}
