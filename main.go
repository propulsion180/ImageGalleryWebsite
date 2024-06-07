package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ImageMeta struct {
	FilePath     string
	Description  string
	ISO          string
	ShutterSpeed string
	Aperture     string
	Location     string
}

type User struct {
	Username string
	Password string
	admin    bool
}

// Initialize the database and create the table if it doesn't exist
func initDB(dbName string) (*sql.DB, error) {
	dbExists := true
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		dbExists = false
	}

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	if !dbExists {
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
			return nil, fmt.Errorf("could not create table: %v", err)
		}

		createTableSQL = `CREATE TABLE IF NOT EXISTS users (
			username TEXT PRIMARY KEY,
			password TEXT NOT NULL,
			admin BOOLEAN NOT NULL
		);`

		_, err = db.Exec(createTableSQL)
		if err != nil {
			return nil, fmt.Errorf("could not create users: %v", err)
		}
	}

	return db, nil
}

// AddImageMeta adds a new image metadata entry to the database
func addImageMeta(db *sql.DB, img ImageMeta) error {
	insertSQL := `INSERT INTO image_metadata (filepath, description, iso, shutterspeed, aperture, location) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(insertSQL, img.FilePath, img.Description, img.ISO, img.ShutterSpeed, img.Aperture, img.Location)
	return err
}

func addUser(db *sql.DB, user User) error {
	fmt.Println(user)
	insertSQL := `INSERT INTO users (username, password, admin) VALUES (?, ?, ?)`
	_, err := db.Exec(insertSQL, user.Username, user.Password, user.admin)
	return err
}

// DeleteImageMeta deletes an image metadata entry from the database
func deleteImageMeta(db *sql.DB, filePath string) error {
	deleteSQL := `DELETE FROM image_metadata WHERE filepath = ?`
	_, err := db.Exec(deleteSQL, filePath)
	return err
}

// GetAllImageMeta retrieves all image metadata entries from the database
func getAllImageMeta(db *sql.DB) (map[string][]ImageMeta, error) {
	querySQL := `SELECT filepath, description, iso, shutterspeed, aperture, location FROM image_metadata`
	rows, err := db.Query(querySQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []ImageMeta
	for rows.Next() {
		var img ImageMeta
		err := rows.Scan(&img.FilePath, &img.Description, &img.ISO, &img.ShutterSpeed, &img.Aperture, &img.Location)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	return map[string][]ImageMeta{"data": images}, nil
}

func getUser(db *sql.DB, username string) (User, error) {
	querySQL := `SELECT username, password, admin FROM users WHERE username = ?`
	row := db.QueryRow(querySQL, username)
	var user User
	err := row.Scan(&user.Username, &user.Password, &user.admin)
	return user, err
}

func getImageMeta(db *sql.DB, filePath string) (ImageMeta, error) {
	querySQL := `SELECT filepath, description, iso, shutterspeed, aperture, location FROM image_metadata WHERE filepath = ?`
	row := db.QueryRow(querySQL, filePath)
	var img ImageMeta
	err := row.Scan(&img.FilePath, &img.Description, &img.ISO, &img.ShutterSpeed, &img.Aperture, &img.Location)
	return img, err
}

func deleteCookie(slice []*http.Cookie, name string) []*http.Cookie {
	for i, cookie := range slice {
		if cookie.Value == name {
			slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}

	return slice
}

func containsCookie(slice []*http.Cookie, name string) bool {
	for _, cookie := range slice {
		if cookie.Value == name {
			return true
		}
	}
	return false
}

func checkUser(db *sql.DB, username string, password string) bool {
	querySQL := `SELECT username, password FROM users WHERE username = ? AND password = ?`
	row := db.QueryRow(querySQL, username, password)
	var user User
	err := row.Scan(&user.Username, &user.Password)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func checkUserAdmin(db *sql.DB, username string) (bool, error) {
	querySQL := `SELECT admin FROM users WHERE username = ?`
	row := db.QueryRow(querySQL, username)
	var admin bool
	err := row.Scan(&admin)
	if err != nil {
		return false, err
	}
	return admin, nil
}

func main() {
	var cookies []*http.Cookie

	db, err := initDB("images.db")
	if err != nil {
		fmt.Println("failed to open database.")
	}

	fs := http.FileServer(http.Dir("images"))
	http.Handle("/images/", http.StripPrefix("/images/", fs))

	fs2 := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs2))

	fmt.Println("Starting Server")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFiles("index.html"))
		data, err := getAllImageMeta(db)
		if err != nil {
			fmt.Println("Failed to get all data")
		}
		fmt.Println(data)
		fmt.Println(r.Cookie("session_token"))
		t.Execute(w, data)
	})

	http.HandleFunc("/admin/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("here")
		c, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		val := c.Value
		u := strings.Split(val, "@")[0]

		admin, err := checkUserAdmin(db, u)
		if err != nil {
			fmt.Println("Failed to check if user is admin")
		}
		fmt.Println("admin: ", admin)

		if !admin {
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}

		fmt.Println("here2")
		t := template.Must(template.ParseFiles("admin.html"))
		data, err := getAllImageMeta(db)
		if err != nil {
			fmt.Println("Failed to get all data admin")
		}
		t.Execute(w, data)
	})

	http.HandleFunc("/admin/addimage/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			fmt.Fprintf(w, "Unable to parse form: %v", err)
			return
		}

		description := r.PostFormValue("description")
		iso := r.PostFormValue("iso")
		shutterSpeed := r.PostFormValue("shutter-speed")
		aperture := r.PostFormValue("aperture")
		location := r.PostFormValue("location")

		// Retrieve the file from form
		file, _, err := r.FormFile("image")
		if err != nil {
			fmt.Fprintf(w, "Error retrieving file: %v", err)
			return
		}
		defer file.Close()

		// Read the first 512 bytes to determine the content type
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil {
			fmt.Fprintf(w, "Error reading file: %v", err)
			return
		}

		// Determine the content type of the file
		contentType := http.DetectContentType(buffer)
		ext, err := mime.ExtensionsByType(contentType)
		if err != nil || len(ext) == 0 {
			fmt.Fprintf(w, "Unable to determine file extension: %v", err)
			return
		}

		// Reset the file read pointer to the beginning
		file.Seek(0, io.SeekStart)

		// Create a temporary file with the correct extension
		tempFile, err := os.CreateTemp("images", fmt.Sprintf("upload-*%s", ext[0]))
		if err != nil {
			fmt.Fprintf(w, "Unable to create temp file: %v", err)
			return
		}
		defer tempFile.Close()

		// Copy the uploaded file's content to the temporary file
		_, err = io.Copy(tempFile, file)
		if err != nil {
			fmt.Fprintf(w, "Error saving file: %v", err)
			return
		}

		// Get the file path
		filePath := tempFile.Name()

		fmt.Println(description)
		fmt.Println(iso)
		fmt.Println(shutterSpeed)
		fmt.Println(aperture)
		fmt.Println(filePath)

		err = addImageMeta(db, ImageMeta{FilePath: filePath, Description: description, ISO: iso, ShutterSpeed: shutterSpeed, Aperture: aperture, Location: location})

		tmpl := template.Must(template.ParseFiles("admin.html"))
		tmpl.ExecuteTemplate(w, "simple-list", ImageMeta{FilePath: filePath, Description: description, ISO: iso, ShutterSpeed: shutterSpeed, Aperture: aperture, Location: location})
	})

	http.HandleFunc("/detail", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		path, err = url.QueryUnescape(path)
		fmt.Println(path)
		t := template.Must(template.ParseFiles("detail.html"))
		data, err := getImageMeta(db, path)
		if err != nil {
			fmt.Println("Failed to get all data admin")
		}
		t.Execute(w, data)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			uname := r.Form.Get("username")
			pword := r.Form.Get("password")
			fmt.Println("Outside of checkUser")
			if checkUser(db, uname, pword) {
				c := http.Cookie{
					Name:     "session_token",
					Value:    uname + "@20",
					HttpOnly: true,
					MaxAge:   3600,
				}
				fmt.Println("before cookie added")
				cookies = append(cookies, &c)
				fmt.Println("after to cookies")
				http.SetCookie(w, &c)
				fmt.Println("after set cookie")
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

		} else {
			fmt.Println(" inhere")
			t := template.Must(template.ParseFiles("login.html"))
			t.Execute(w, nil)
		}

	})

	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			fmt.Println("Inside of signup")
			r.ParseForm()
			uname := r.Form.Get("username")
			pword := r.Form.Get("password")
			admin := r.Form.Get("admin")

			fmt.Println(uname, pword, admin)
			ts := checkUser(db, uname, pword)
			fmt.Println(ts)
			if !ts {
				fmt.Println("User not already exists")
				addUser(db, User{Username: uname, Password: pword, admin: admin == "on"})
				http.Redirect(w, r, "/login", http.StatusFound)
			}

		} else {
			t := template.Must(template.ParseFiles("signup.html"))
			t.Execute(w, nil)
		}
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			fmt.Println("Failed to get cookie")
			http.Redirect(w, r, "/", http.StatusFound)
		}

		c.MaxAge = -1
		c.Expires = time.Unix(0, 0)

		http.SetCookie(w, c)
		deleteCookie(cookies, c.Value)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	log.Fatal(http.ListenAndServe(":9000", nil))
}
