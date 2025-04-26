package handlers

import (
	"encoding/json"
	"gallery-server/auth"
	"gallery-server/db"
	"gallery-server/models"
	"log"
	"net/http"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var login_data models.LoginData
	err := json.NewDecoder(r.Body).Decode(&login_data)
	if err != nil {
		log.Println("failed to decode json data for login: ", err.Error())
		http.Error(w, "Invalid Request payload", http.StatusBadRequest)
		return
	}

	conn := db.ConnectDB()
	defer conn.Close()

	passCorrect, err := db.VerifyPassword(conn, login_data.Username, login_data.Password)
	user, err := db.GetUser(conn, login_data.Username, login_data.Password)

	if err != nil {
		log.Println("caught error from getting user: ", err.Error())
		http.Error(w, "Error while getting user", http.StatusInternalServerError)
		return
	}

	if user == nil {
		log.Println("invalid password sent by " + login_data.Username)
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	tkn, err := auth.GenerateJWT(user.Username, user.Admin)

	if err != nil {
		log.Println("caught error from generate JWT: ", err.Error())
		http.Error(w, "Error while generating token", http.StatusInternalServerError)
		return
	}

	res := db.SetToken(conn, tkn, user)

	if !res {
		log.Println("failed to set token to user")
		http.Error(w, "failed to set user's token", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    tkn,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	response := models.LoginResponse{
		Message:  "Sucessful login",
		Username: user.Username,
		Admin:    user.Admin,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
