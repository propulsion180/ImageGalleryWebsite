package handlers

import (
	"encoding/json"
	"gallery-server/auth"
	"gallery-server/db"
	"gallery-server/models"
	"log"
	"net/http"
)

func AllImageHandler(w http.ResponseWriter, r *http.Request) {
	conn := db.ConnectDB()
	defer conn.Close()
	images, err := db.GetAllImageMeta(conn)
	if err != nil {
		log.Println("caught error from getall image meta: ", err.Error())
		http.Error(w, "failed to get all images", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(images); err != nil {
		log.Println("failed to encode the images slice: ", err.Error())
		http.Error(w, "failed to encode images", http.StatusInternalServerError)
		return
	}
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	var image_data models.SingleImageData
	err := json.NewDecoder(r.Body).Decode(&image_data)
	if err != nil {
		log.Println("failed to decode the json data for single image data: ", err.Error())
		http.Error(w, "Invalid Request payload", http.StatusBadRequest)
		return
	}
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		log.Println("failed to get cookie from requerst: ", err.Error())
		http.Error(w, "Can't get cookie for verification", http.StatusBadRequest)
		return
	}
	_, err = auth.VerifyJWT(cookie.Value)
	if err != nil {
		log.Println("unauthorized user tried to get an image or failed to verify: ", err.Error())
		http.Error(w, "unauthorized user or failed token verification", http.StatusBadRequest)
		return
	}
	conn := db.ConnectDB()
	defer conn.Close()
	image, err := db.GetImageMeta(conn, image_data.FilePath)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(image); err != nil {
		log.Println("failed to encode image: ", err.Error())
		http.Error(w, "failed to encode image", http.StatusInternalServerError)
		return
	}
}
