package handlers

import (
	"encoding/json"
	"fmt"
	"gallery-server/auth"
	"gallery-server/db"
	"gallery-server/models"
	"image"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

func AllImageHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("All Images have been requested")
	conn := db.ConnectDB()
	defer conn.Close()
	images, err := db.GetAllImageMeta(conn)
	if err != nil {
		slog.Error("Failed to get all images", "status", http.StatusInternalServerError, "error", err.Error())
		http.Error(w, "failed to get all images", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(images); err != nil {
		slog.Error("Failed to encode the images slice", "status", http.StatusInternalServerError, "error", err.Error())
		http.Error(w, "failed to encode images", http.StatusInternalServerError)
		return
	}
	slog.Info("All images sent!", "status", http.StatusOK)
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	var image_data models.SingleImageData
	err := json.NewDecoder(r.Body).Decode(&image_data)
	if err != nil {
		slog.Error("Failed to decode the JSON data in single image data", "status", http.StatusBadRequest, "error", err.Error())
		http.Error(w, "Invalid Request payload", http.StatusBadRequest)
		return
	}
	slog.Info("Image requested", "filepath", image_data.FilePath)
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		slog.Error("Failed to get cookie from request: ", "status", http.StatusBadRequest, "error", err.Error())
		http.Error(w, "Can't get cookie for verification", http.StatusBadRequest)
		return
	}
	_, err = auth.VerifyJWT(cookie.Value)
	if err != nil {
		slog.Error("Unauthorized user tried to get an image or failed to verify ", "status", http.StatusUnauthorized, "error", err.Error())
		http.Error(w, "unauthorized user or failed token verification", http.StatusUnauthorized)
		return
	}
	conn := db.ConnectDB()
	defer conn.Close()
	image, err := db.GetImageMeta(conn, image_data.FilePath)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(image); err != nil {
		slog.Error("Failed to encode image: ", "status", http.StatusInternalServerError, "error", err.Error())
		http.Error(w, "failed to encode image", http.StatusInternalServerError)
		return
	}
	slog.Info("Image sent to user", "status", http.StatusOK)
}

func applyOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 3:
		return imaging.Rotate180(img)
	case 6:
		return imaging.Rotate270(img)
	case 8:
		return imaging.Rotate90(img)
	default:
		return img
	}
}

func NewImageHandler(w http.ResponseWriter, r *http.Request) {

	var img models.ImageMeta
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		slog.Error("Failed to parse form data: ", "status", http.StatusBadRequest, "error", err.Error())
		http.Error(w, "Error parsing the form", http.StatusBadRequest)
		return
	}

	img.Description = r.FormValue("description")
	img.ShutterSpeed = r.FormValue("shutterspeed")
	img.ISO = r.FormValue("iso")
	img.Aperture = r.FormValue("aperture")
	img.Location = r.FormValue("location")
	file, handler, err := r.FormFile("file")

	if err != nil {
		slog.Error("Failed to retrieive the file ", "status", http.StatusInternalServerError, "error", err.Error())
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	outputDir := "images"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			slog.Error("Failed to create images directory", "status", http.StatusInternalServerError, "error", err.Error())
			http.Error(w, "Failed to create images directory", http.StatusInternalServerError)
			return
		}
	}

	image, _, err := image.Decode(file)
	if err != nil {
		slog.Error("Failed to decode image when saving an image", "status", http.StatusInternalServerError, "error", err.Error())
		http.Error(w, "failed to decode image", http.StatusInternalServerError)
		return
	}

	file.Seek(0, io.SeekStart)

	x, err := exif.Decode(file)

	if err != nil {
		slog.Error("Failed to decode image exif data", "status", http.StatusInternalServerError, "error", err.Error())
		http.Error(w, "failed to decode image exif data", http.StatusInternalServerError)
		return
	}

	orientation := 1

	tag, err := x.Get(exif.Orientation)

	if err != nil {
		slog.Error("Failed to get exif orientation", "status", http.StatusInternalServerError, "error", err.Error())
		http.Error(w, "failed to get exif orientation", http.StatusInternalServerError)
		return
	}

	orientation, _ = tag.Int(0)

	base := img.FilePath
	if base == "" {
		base = handler.Filename
	}
	base = strings.TrimSuffix(base, filepath.Ext(base))
	newFileName := base + ".webp"
	outputPath := filepath.Join(outputDir, newFileName)

	outFile, err := os.Create(outputPath)
	if err != nil {
		slog.Error("Failed to create output file path", "status", http.StatusInternalServerError, "error", err.Error())
		http.Error(w, "failed to create image output path", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()
	img.FilePath = outputPath

	rImage := applyOrientation(image, orientation)

	options := &webp.Options{Quality: 80}
	if err := webp.Encode(outFile, rImage, options); err != nil {
		slog.Error("Failed to encode image into webp", "status", http.StatusInternalServerError, "error", err.Error())
		http.Error(w, "failed to encode image into webp", http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		slog.Error("Failed to get cookie from request", "status", http.StatusBadGateway, "error", err.Error())
		http.Error(w, "Can't get cookie for verification", http.StatusBadRequest)
		return
	}
	claims, err := auth.VerifyJWT(cookie.Value)
	if err != nil {
		slog.Error("Unauthorized user tried to get add image or failed to verify", "status", http.StatusBadRequest, "error", err.Error())
		http.Error(w, "unauthorized user or failed token verification", http.StatusBadRequest)
		return
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		slog.Error("Failed to get sub from claims", "status", http.StatusInternalServerError)
		http.Error(w, "failed to get sub from cookie claims", http.StatusInternalServerError)
		return
	}
	auth, ok := claims["admin"].(bool)
	if !ok {
		slog.Error("Failed to get auth from claims", "status", http.StatusUnauthorized)
		http.Error(w, "failed to auth user", http.StatusUnauthorized)
		return
	}

	if !auth {
		slog.Error("Not auhorized to add image", "status", http.StatusUnauthorized)
		http.Error(w, "user not auth", http.StatusUnauthorized)
		return
	}

	conn := db.ConnectDB()
	defer conn.Close()
	fmt.Println(img)
	err = db.AddImageMeta(conn, &img, sub)
	if err != nil {
		slog.Error("Failed to add the image entry to the database", "status", http.StatusInternalServerError, "error", err.Error())
		http.Error(w, "failed to add image entry to the database: ", http.StatusInternalServerError)
		return
	}

	slog.Info("Image sucessfully uploaded and converted", "status", http.StatusOK, "filepath", outputPath)

	w.WriteHeader(http.StatusOK)
}

func UpdateImageHandler(w http.ResponseWriter, r *http.Request) {
	var img models.ImageMeta
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		slog.Error("Failed to parse form data: ", "status", http.StatusBadRequest, "error", err.Error())
		http.Error(w, "Error parsing the form", http.StatusBadRequest)
		return
	}

	img.FilePath = r.FormValue("filepath")
	img.Description = r.FormValue("description")
	img.ShutterSpeed = r.FormValue("shutterspeed")
	img.ISO = r.FormValue("iso")
	img.Aperture = r.FormValue("aperture")
	img.Location = r.FormValue("location")

	slog.Info("Updating Image", "filepath", img.FilePath)

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		slog.Error("Failed to get cookie from request: ", "status", http.StatusBadRequest, "error", err.Error())
		http.Error(w, "Can't get cookie for verification", http.StatusBadRequest)
		return
	}
	claims, err := auth.VerifyJWT(cookie.Value)
	if err != nil {
		slog.Error("Unauthorized user tried to get add image or failed to verify: ", "status", http.StatusBadRequest, "error", err.Error())
		http.Error(w, "unauthorized user or failed token verification", http.StatusBadRequest)
		return
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		slog.Error("Failed to get sub from claims", "status", http.StatusInternalServerError)
		http.Error(w, "failed to get sub from cookie claims", http.StatusInternalServerError)
		return
	}
	auth, ok := claims["admin"].(bool)
	if !ok || !auth {
		slog.Error("Failed to get auth from claims", "status", http.StatusUnauthorized)
		http.Error(w, "failed to auth user", http.StatusUnauthorized)
		return
	}

	conn := db.ConnectDB()
	defer conn.Close()

	err = db.SetImageMeta(conn, &img, sub)
	if err != nil {
		slog.Error("Failed to update image entry to the database:", "status", http.StatusInternalServerError, "error", err.Error())
		http.Error(w, "failed to add image entry to the database: ", http.StatusInternalServerError)
		return
	}

	slog.Info("Image sucessfully updated image entry", "status", http.StatusOK)

	w.WriteHeader(http.StatusOK)
}

func DeleteImageHandler(w http.ResponseWriter, r *http.Request) {
	var image_data models.SingleImageData
	err := json.NewDecoder(r.Body).Decode(&image_data)

	slog.Info("Deleting Image", "filepaht", image_data.FilePath)

	if err != nil {
		slog.Error("Failed to decode the json data for single image data", "status", http.StatusBadRequest, "error", err.Error())
		http.Error(w, "Invalid Request payload", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		slog.Error("Failed to get cookie from request", "status", http.StatusBadRequest, "error", err.Error())
		http.Error(w, "Can't get cookie for verification", http.StatusBadRequest)
		return
	}
	claims, err := auth.VerifyJWT(cookie.Value)
	if err != nil {
		slog.Error("Unauthorized user tried to get add image or failed to verify", "status", http.StatusBadRequest, "error", err.Error())
		http.Error(w, "unauthorized user or failed token verification", http.StatusBadRequest)
		return
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		slog.Error("Failed to get sub from claims", "status", http.StatusInternalServerError)
		http.Error(w, "failed to get sub from cookie claims", http.StatusInternalServerError)
		return
	}
	auth, ok := claims["admin"].(bool)
	if !ok || !auth {
		slog.Error("Fialed to get auth from claims", "status", http.StatusUnauthorized)
		http.Error(w, "failed to auth user", http.StatusUnauthorized)
		return
	}

	conn := db.ConnectDB()
	defer conn.Close()

	err = db.DeleteImageMeta(conn, image_data.FilePath, sub)
	if err != nil {
		slog.Error("Failed to delete the image entry from the database", "status", http.StatusInternalServerError, "error", err.Error())
		http.Error(w, "failed to delete the image entry", http.StatusInternalServerError)
		return
	}

	slog.Info("Image sucessfully deleted", "status", http.StatusOK)

	w.WriteHeader(http.StatusOK)
}
