package handlers

import (
	"gallery-server/auth"
	"gallery-server/db"
	"log"
	"net/http"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookieToBeInvalidated, err := r.Cookie("auth_token")
	expired := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, expired)
	conn := db.ConnectDB()
	defer conn.Close()
	claims, err := auth.VerifyJWT(cookieToBeInvalidated.Value)
	if err != nil {
		log.Println("failed to verify token and get claims: ", err.Error())
		http.Error(w, "failed to verify token", http.StatusInternalServerError)
		return
	}
	subClaim, ok := claims["sub"].(string)
	if !ok {
		log.Println("failed to get the sub from the token's claims")
		http.Error(w, "Failed to parse cookie token", http.StatusInternalServerError)
		return
	}
	err = db.DeleteToken(conn, subClaim)
	if err != nil {
		log.Println("caught error from delete token", err.Error())
		http.Error(w, "failed to delete token from database", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
