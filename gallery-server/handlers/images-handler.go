package handlers

import (
	"encoding/json"
	"gallery-server/auth"
	"gallery-server/db"
	"gallery-server/models"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"

	_ "image/jpeg"
	_ "image/png"
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

func NewImageHandler(w http.ResponseWriter, r *http.Request) {

	var img models.ImageMeta
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		log.Println("failed to parse form data: ", err.Error())
		http.Error(w, "Error parsing the form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		log.Println("failed to retrieive the file: ", err.Error())
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	img.FilePath = r.FormValue("filename")
	img.Description = r.FormValue("description")
	img.ShutterSpeed = r.FormValue("shutterspeed")
	img.ISO = r.FormValue("iso")
	img.Aperture = r.FormValue("aperture")
	img.Location = r.FormValue("location")

	outputDir := filepath.Join("..", "images")
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			log.Println("failed to create images directory: ", err.Error())
			http.Error(w, "Failed to create images directory", http.StatusInternalServerError)
			return
		}
	}

	image, _, err := image.Decode(file)
	if err != nil {
		log.Println("failde to decode image when saving an image: ", err.Error())
		http.Error(w, "failed to decode image", http.StatusInternalServerError)
		return
	}

	base := img.FilePath
	if base == "" {
		base = handler.Filename
	}
	base = strings.TrimSuffix(base, filepath.Ext(base))
	newFileName := base + ".webp"
	outputPath := filepath.Join(outputDir, newFileName)

	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Println("failed to create output file path: ", err.Error())
		http.Error(w, "failed to create image output path", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	img.FilePath = outputPath

	options := &webp.Options{Quality: 80}
	if err := webp.Encode(outFile, image, options); err != nil {
		log.Println("failed to encode image into webp: ", err.Error())
		http.Error(w, "failed to encode image into webp", http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		log.Println("failed to get cookie from request: ", err.Error())
		http.Error(w, "Can't get cookie for verification", http.StatusBadRequest)
		return
	}
	claims, err := auth.VerifyJWT(cookie.Value)
	if err != nil {
		log.Println("unauthorized user tried to get add image or failed to verify: ", err.Error())
		http.Error(w, "unauthorized user or failed token verification", http.StatusBadRequest)
		return
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		log.Println("failed to get sub from claims")
		http.Error(w, "failed to get sub from cookie claims", http.StatusInternalServerError)
		return
	}
	auth, ok := claims["auth"].(bool)
	if !ok || !auth {
		log.Println("fialed to get auth from claims")
		http.Error(w, "failed to auth user", http.StatusUnauthorized)
		return
	}

	conn := db.ConnectDB()
	defer conn.Close()

	err = db.AddImageMeta(conn, &img, sub)
	if err != nil {
		log.Println("failed to add the image entry to the database: ", err.Error())
		http.Error(w, "failed to add image entry to the database: ", http.StatusInternalServerError)
		return
	}

	log.Println("Image sucessfully uploaded and converted: ", outputPath)

	w.WriteHeader(http.StatusOK)
}

func UpdateImageHandler(w http.ResponseWriter, r *http.Request) {
	var img models.ImageMeta
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		log.Println("failed to parse form data: ", err.Error())
		http.Error(w, "Error parsing the form", http.StatusBadRequest)
		return
	}

	img.FilePath = r.FormValue("filename")
	img.Description = r.FormValue("description")
	img.ShutterSpeed = r.FormValue("shutterspeed")
	img.ISO = r.FormValue("iso")
	img.Aperture = r.FormValue("aperture")
	img.Location = r.FormValue("location")

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		log.Println("failed to get cookie from request: ", err.Error())
		http.Error(w, "Can't get cookie for verification", http.StatusBadRequest)
		return
	}
	claims, err := auth.VerifyJWT(cookie.Value)
	if err != nil {
		log.Println("unauthorized user tried to get add image or failed to verify: ", err.Error())
		http.Error(w, "unauthorized user or failed token verification", http.StatusBadRequest)
		return
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		log.Println("failed to get sub from claims")
		http.Error(w, "failed to get sub from cookie claims", http.StatusInternalServerError)
		return
	}
	auth, ok := claims["auth"].(bool)
	if !ok || !auth {
		log.Println("fialed to get auth from claims")
		http.Error(w, "failed to auth user", http.StatusUnauthorized)
		return
	}

	conn := db.ConnectDB()
	defer conn.Close()

	err = db.SetImageMeta(conn, &img, sub)
	if err != nil {
		log.Println("failed to update image entry to the database:", err.Error())
		http.Error(w, "failed to add image entry to the database: ", http.StatusInternalServerError)
		return
	}

	log.Println("Image sucessfully updated image entry")

	w.WriteHeader(http.StatusOK)
}

func DeleteImageHandler(w http.ResponseWriter, r *http.Request) {
	var image_data models.SingleImageData
	err := json.NewDecoder(r.Body).Decode(&image_data)
	if err != nil {
		log.Println("failed to decode the json data for single image data: ", err.Error())
		http.Error(w, "Invalid Request payload", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		log.Println("failed to get cookie from request: ", err.Error())
		http.Error(w, "Can't get cookie for verification", http.StatusBadRequest)
		return
	}
	claims, err := auth.VerifyJWT(cookie.Value)
	if err != nil {
		log.Println("unauthorized user tried to get add image or failed to verify: ", err.Error())
		http.Error(w, "unauthorized user or failed token verification", http.StatusBadRequest)
		return
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		log.Println("failed to get sub from claims")
		http.Error(w, "failed to get sub from cookie claims", http.StatusInternalServerError)
		return
	}
	auth, ok := claims["auth"].(bool)
	if !ok || !auth {
		log.Println("fialed to get auth from claims")
		http.Error(w, "failed to auth user", http.StatusUnauthorized)
		return
	}

	conn := db.ConnectDB()
	defer conn.Close()

	err = db.DeleteImageMeta(conn, image_data.FilePath, sub)
	if err != nil {
		log.Println("failed to delete the image entry from the database: ", err.Error())
		http.Error(w, "failed to delete the image entry", http.StatusInternalServerError)
		return
	}

	log.Println("Image sucessfully deleted")

	w.WriteHeader(http.StatusOK)
}
