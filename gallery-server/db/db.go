package db

import (
	"database/sql"
	"errors"
	"gallery-server/auth"
	"gallery-server/models"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DBName string

func InitDB(dbName string) bool {
	if _, err := os.Stat(dbName); os.IsExist(err) {
		slog.Error("Database initialiation not needed file already exists", "error", err.Error())
		DBName = dbName
		return true
	} else {
		db, err := sql.Open("sqlite3", dbName)
		if err != nil {
			slog.Error("Database initialization failed to open database", "error", err.Error())
			return false
		}

		createTableSQL := `CREATE TABLE IF NOT EXISTS image_metadata (
			filepath TEXT PRIMARY KEY,
			description TEXT,
			iso TEXT,
			shutterspeed TEXT,
			aperture TEXT,
			location TEXT
		);`

		_, err = db.Exec(createTableSQL)
		if err != nil {
			slog.Error("Failed to create image_metadata table in initialization.", "error", err.Error())
			return false
		}

		createTableSQL = `CREATE TABLE IF NOT EXISTS users (
			username TEXT PRIMARY KEY,
			password TEXT NOT NULL,
			admin BOOLEAN NOT NULL,
			token TEXT NOT NULL DEFAULT ''
		);`

		_, err = db.Exec(createTableSQL)
		if err != nil {
			slog.Error("Failed to create the users table in initialization.", "error", err.Error())
			return false
		}
		DBName = dbName
		return true
	}
}

func ConnectDB() *sql.DB {
	if _, err := os.Stat(DBName); os.IsNotExist(err) {
		slog.Error("Failed to connect to database as it dosen't exist in the file system.", "error", err.Error())
		return nil
	}

	db, err := sql.Open("sqlite3", DBName)
	if err != nil {
		slog.Error("Failed to open database even though it exists in file system.", "error", err.Error())
		return nil
	}
	return db
}

//token functions

func SetToken(db *sql.DB, tkn string, usr *models.User) bool {
	stmt := `UPDATE users SET token=? WHERE username=?`
	_, err := db.Exec(stmt, tkn, usr.Username)
	if err != nil {
		slog.Error("Failed to update users with the tokens", "error", err.Error())
		return false
	}
	return true
}

func CheckToken(db *sql.DB, tkn string, usr models.User) bool {
	stmt := `SELECT token FROM users WHERE username=?;`
	row := db.QueryRow(stmt, usr.Username)
	var token string
	err := row.Scan(&token)
	if err != nil {
		slog.Error("Failed to scan row when checking token.", "error", err.Error())
	}
	return token == tkn
}

func DeleteToken(db *sql.DB, username string) error {
	stmt := `UPDATE users SET token=? WHERE username=?;`
	_, err := db.Exec(stmt, "", username)

	if err != nil {
		slog.Error("Failed to delete the token the error was this", "error", err.Error())
		return err
	}
	return nil
}

//user functions

func AddUser(db *sql.DB, user *models.User) (bool, error) {
	if !UserCheck(user) {
		slog.Warn("Fails users data checks")
		return false, nil
	}

	hashed, err := auth.HashPassword(user.Password)
	if err != nil {
		return false, err
	}

	statement := `INSERT INTO users (username, password, admin) VALUES (?, ?, ?)`
	_, err = db.Exec(statement, user.Username, hashed, user.Admin)
	if err != nil {
		slog.Error("Failed to add user here is the error", "error", err.Error())
		return false, err
	}
	return true, nil
}

func IsAdmin(db *sql.DB, unameOrToken string) (bool, error) {
	statement := `SELECT admin FROM users WHERE username=? or token=?;`
	row := db.QueryRow(statement, unameOrToken, unameOrToken)
	var admn bool
	err := row.Scan(&admn)
	if err != nil {
		slog.Error("Failed to scan row when checking admin", err.Error())
		return false, err
	}
	if admn {
		slog.Warn("Admin Logged In!", "unameortoken", unameOrToken)
	}
	return admn, nil
}

func VerifyPassword(db *sql.DB, username string, password string) (bool, error) {
	stmt := `SELECT password FROM users WHERE username=?;`
	row := db.QueryRow(stmt, username)
	var pword string
	err := row.Scan(&pword)
	if err != nil {
		slog.Error("Failed to get password of user to verify", "username", username, "password", password, "error", err.Error())
		return false, err
	}
	return auth.ValidateHash(password, pword), nil
}

func GetUser(db *sql.DB, uname string, pword string) (*models.User, error) {
	veri, err := VerifyPassword(db, uname, pword)
	if err != nil {
		slog.Error("Error during the verification of password for getuser", "error", err.Error())
		return nil, err
	}
	if !veri {
		slog.Error("Verification failed in getuser", "username", uname)
		return nil, nil
	}

	is_admin, err := IsAdmin(db, uname)

	if err != nil {
		slog.Error("Error during admin check in getuser", "error", err.Error())
		return nil, err
	}

	user := models.User{Username: uname, Password: pword, Admin: is_admin}

	return &user, nil
}

// AddImageMeta adds a new image metadata entry to the database
func AddImageMeta(db *sql.DB, img *models.ImageMeta, adder string) error {

	if !CameraCheck(img) {
		slog.Warn("Image data fails checks")
		return errors.New("image data fails checks")
	}

	ad, err := IsAdmin(db, adder)
	if err != nil {
		slog.Error("Something went wrong when checking admin priveleges when adding image", "error", err.Error())
		return err
	}
	if !ad {
		slog.Warn("Unprivleged user tried to add an image", "username", adder, "filepath", img.FilePath)
		return nil
	}
	insertSQL := `INSERT INTO image_metadata (filepath, description, iso, shutterspeed, aperture, location) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(insertSQL, img.FilePath, img.Description, img.ISO, img.ShutterSpeed, img.Aperture, img.Location)
	if err != nil {
		slog.Error("Failed to add image to database:", "error", err.Error())
		return err
	}
	return nil
}

// DeleteImageMeta deletes an image metadata entry from the database
func DeleteImageMeta(db *sql.DB, filePath string, deleter string) error {
	ad, err := IsAdmin(db, deleter)
	if err != nil {
		slog.Error("Failed to check admin priveleges when deleting image", "error", err.Error())
		return err
	}
	if !ad {
		slog.Warn("Unpriveleged user tried to delete an image", "username", deleter, "filepath", filePath)
		return nil
	}
	deleteSQL := `DELETE FROM image_metadata WHERE filepath = ?`
	_, err = db.Exec(deleteSQL, filePath)
	if err != nil {
		slog.Error("Failed to delete image from database", "error", err.Error())
		return err
	}
	return nil
}

func SetImageMeta(db *sql.DB, img *models.ImageMeta, setter string) error {
	if !CameraCheck(img) {
		slog.Warn("Image data fails checks")
		return errors.New("image data fails checks")
	}
	ad, err := IsAdmin(db, setter)
	if err != nil {
		slog.Error("Failed to check admin when setting image", "error", err.Error())
		return err
	}
	if !ad {
		slog.Warn("Unprivileged user tried to update an image", "username", setter, "filepath", img.FilePath)
		return nil
	}
	updateSQL := `UPDATE image_metadata SET description = ?, iso = ?, shutterspeed = ?, aperture = ?, location = ? WHERE filepath = ?`
	_, err = db.Exec(updateSQL, img.Description, img.ISO, img.ShutterSpeed, img.Aperture, img.Location, img.FilePath)
	if err != nil {
		slog.Error("Failed to set image properties", "error", err.Error())
		return err
	}
	return nil
}

// GetAllImageMeta retrieves all image metadata entries from the database
func GetAllImageMeta(db *sql.DB) ([]models.ImageMeta, error) {
	querySQL := `SELECT filepath, description, iso, shutterspeed, aperture, location FROM image_metadata`
	rows, err := db.Query(querySQL)
	if err != nil {
		slog.Error("Faailed to get all images from the database", "error", err.Error())
		return nil, err
	}
	defer rows.Close()

	var images []models.ImageMeta
	for rows.Next() {
		var img models.ImageMeta
		err := rows.Scan(&img.FilePath, &img.Description, &img.ISO, &img.ShutterSpeed, &img.Aperture, &img.Location)
		if err != nil {
			slog.Error("Failed to scan the rows into ImageMeta structs", "error", err.Error())
			return nil, err
		}
		images = append(images, img)
	}

	return images, nil
}

func GetImageMeta(db *sql.DB, filePath string) (*models.ImageMeta, error) {
	querySQL := `SELECT filepath, description, iso, shutterspeed, aperture, location FROM image_metadata WHERE filepath = ?`
	row := db.QueryRow(querySQL, filePath)
	var img models.ImageMeta
	err := row.Scan(&img.FilePath, &img.Description, &img.ISO, &img.ShutterSpeed, &img.Aperture, &img.Location)
	if err != nil {
		slog.Error("Failed to scan image meta from the databse", "error", err.Error())
		return nil, err
	}
	return &img, nil
}
