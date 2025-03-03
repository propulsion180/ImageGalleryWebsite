package db

import (
	"database/sql"
	"gallery-server/models"
	"log"
	"os"
)

var DBName string

func InitDB(dbName string) bool {
	if _, err := os.Stat(dbName); os.IsExist(err) {
		log.Println("initialiation not needed file already exists")
		DBName = dbName
		return true
	} else {
		db, err := sql.Open("sqlite3", dbName)
		if err != nil {
			log.Println("initialization failed to open database. Error is: ", err.Error())
			return false
		}

		createTableSQL := `CREATE TABLE IF NOT EXISTS image_metadata (
			filepath TEXT PRIMARY KEY,
			description TEXT,
			iso TEXT,
			shutterspeed TEXT,
			aperture TEXT,
			location TEXT,
		);`

		_, err = db.Exec(createTableSQL)
		if err != nil {
			log.Println("Failed to create image_metadata table in initialization. Error is: ", err.Error())
			return false
		}

		createTableSQL = `CREATE TABLE IF NOT EXISTS users (
			username TEXT PRIMARY KEY,
			password TEXT NOT NULL,
			admin BOOLEAN NOT NULL,
			token TEXT NOT NULL DEFAULT 
		);`

		_, err = db.Exec(createTableSQL)
		if err != nil {
			log.Println("Failed to create the users table in initialization. Error is:", err.Error())
			return false
		}
		DBName = dbName
		return true
	}
}

func ConnectDB() *sql.DB {
	if _, err := os.Stat(DBName); os.IsNotExist(err) {
		log.Println("failed to connect to database as it dosen't exist in the file system.")
		return nil
	}

	db, err := sql.Open("sqlite3", DBName)
	if err != nil {
		log.Println("failed to open database even though it exists in file system. Here is the error: ", err.Error())
		return nil
	}
	return db
}

//token functions

func SetToken(db *sql.DB, tkn string, usr *models.User) bool {
	stmt := `UPDATE users SET token=? WHERE username=?`
	_, err := db.Exec(stmt, tkn, usr.Username)
	if err != nil {
		log.Println("failed to update users with the token, here is the error: ", err.Error())
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
		log.Println("failed to scan row when checking token. here is the error: ", err.Error())
	}
	return token == tkn
}

func DeleteToken(db *sql.DB, username string) error {
	stmt := `UPDATE users SET token=? WHERE username=?;`
	_, err := db.Exec(stmt, "", username)

	if err != nil {
		log.Println("failed to delete the token the error was this: ", err.Error())
		return err
	}
	return nil
}

//user functions

func AddUser(db *sql.DB, user models.User, adder string) (bool, error) {
	adminAdder, err := IsAdmin(db, adder)
	if err != nil {
		log.Println("failed to check admin privileges when adding a user")
		return false, err
	}
	if !adminAdder {
		log.Println("unpriviledged user " + adder + "tried to add a new user")
		return false, nil
	}

	statement := `INSERT INTO users (username, password, admin) VALUES (?, ?, ?)`
	_, err = db.Exec(statement, user.Username, user.Password, user.Admin)
	if err != nil {
		log.Println("Failed to add user here is the error: ", err.Error())
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
		log.Println("Failed to scan row: ", err.Error())
		return false, err
	}
	return admn, nil
}

func VerifyPassword(db *sql.DB, username string, password string) (bool, error) {
	stmt := `SELECT password FROM users WHERE username=?;`
	row := db.QueryRow(stmt, username)
	var pword string
	err := row.Scan(&pword)
	if err != nil {
		log.Println("failed to get password of user to verify, here is the error: ", err.Error())
		return false, err
	}
	return pword == password, nil
}

func GetUser(db *sql.DB, uname string, pword string) (*models.User, error) {
	veri, err := VerifyPassword(db, uname, pword)
	if err != nil {
		log.Println("error during the verification of password for getuser")
		return nil, err
	}
	if !veri {
		log.Println("verification in getuser. either password or username is wrong")
		return nil, nil
	}

	is_admin, err := IsAdmin(db, uname)

	if err != nil {
		log.Println("error during admin check in getuser")
		return nil, err
	}

	user := models.User{Username: uname, Password: pword, Admin: is_admin}

	return &user, nil
}

// AddImageMeta adds a new image metadata entry to the database
func AddImageMeta(db *sql.DB, img *models.ImageMeta, adder string) error {
	ad, err := IsAdmin(db, adder)
	if err != nil {
		log.Fatal("something went wrong when checking admin priveleges when adding image:", err.Error())
		return err
	}
	if !ad {
		log.Println("unprivleged user " + adder + "tried to add an image")
		return nil
	}
	insertSQL := `INSERT INTO image_metadata (filepath, description, iso, shutterspeed, aperture, location) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(insertSQL, img.FilePath, img.Description, img.ISO, img.ShutterSpeed, img.Aperture, img.Location)
	if err != nil {
		log.Println("failed to add image to database: ", err.Error())
		return err
	}
	return nil
}

// DeleteImageMeta deletes an image metadata entry from the database
func DeleteImageMeta(db *sql.DB, filePath string, deleter string) error {
	ad, err := IsAdmin(db, deleter)
	if err != nil {
		log.Println("failed to check admin priveleges when deleting image: ", err.Error())
		return err
	}
	if !ad {
		log.Println("unprivelegede user " + deleter + "tried to delete and image")
		return nil
	}
	deleteSQL := `DELETE FROM image_metadata WHERE filepath = ?`
	_, err = db.Exec(deleteSQL, filePath)
	if err != nil {
		log.Println("failed to delete image from database")
		return err
	}
	return nil
}

func SetImageMeta(db *sql.DB, img *models.ImageMeta, setter string) error {
	ad, err := IsAdmin(db, setter)
	if err != nil {
		log.Println("failed to check admin when setting image: ", err.Error())
		return err
	}
	if !ad {
		log.Println("unprivileged user " + setter + "tried to set an image")
		return nil
	}
	updateSQL := `UPDATE image_metadata SET description = ?, iso = ?, shutterspeed = ?, aperture = ?, location = ? WHERE filepath = ?`
	_, err = db.Exec(updateSQL, img.Description, img.ISO, img.ShutterSpeed, img.Aperture, img.Location, img.FilePath)
	if err != nil {
		log.Println("failed to set image properties: ", err.Error())
		return err
	}
	return nil
}

// GetAllImageMeta retrieves all image metadata entries from the database
func GetAllImageMeta(db *sql.DB) ([]models.ImageMeta, error) {
	querySQL := `SELECT filepath, description, iso, shutterspeed, aperture, location FROM image_metadata`
	rows, err := db.Query(querySQL)
	if err != nil {
		log.Println("failed to get all images from the database: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	var images []models.ImageMeta
	for rows.Next() {
		var img models.ImageMeta
		err := rows.Scan(&img.FilePath, &img.Description, &img.ISO, &img.ShutterSpeed, &img.Aperture, &img.Location)
		if err != nil {
			log.Println("failed to scan the rows into ImageMeta structs: ", err.Error())
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
		log.Println("failed to scan image meta from the databse: ", err.Error())
		return nil, err
	}
	return &img, nil
}
