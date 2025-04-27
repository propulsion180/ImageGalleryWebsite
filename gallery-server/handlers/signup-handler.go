package handlers

import (
	"encoding/json"
	"gallery-server/db"
	"gallery-server/models"
	"log"
	"net/http"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var signup_data models.LoginData // Reused for signup
	err := json.NewDecoder(r.Body).Decode(&signup_data)
	if err != nil {
		log.Println("failed to decode json data for signup: ", err.Error())
		http.Error(w, "Invalid Request payload", http.StatusBadRequest)
		return
	}

	conn := db.ConnectDB()
	defer conn.Close()

	worked, err := db.AddUser(conn, &models.User{Username: signup_data.Username, Password: signup_data.Password, Admin: false})

	if !worked {
		log.Println("failed to add user")
		http.Error(w, "Failed to add user", http.StatusUnauthorized)
		return
	}

	if err != nil {
		log.Println("Failed to add user:", err.Error())
		http.Error(w, "Error caught while add user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
