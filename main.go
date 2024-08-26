package main

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"image"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rwcarlsen/goexif/exif"
)

func convertToSHA256(in string) string {
	hash := sha256.New()
	hash.Write([]byte(in))
	checksum := hash.Sum(nil)
	return hex.EncodeToString(checksum)
}

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

func printIP(r *http.Request) {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IPs; the first one is the client IP
		ip := strings.Split(forwarded, ",")[0]
		fmt.Println(strings.TrimSpace(ip))
	}

	// Check the X-Real-Ip header for proxies/load balancers
	realIP := r.Header.Get("X-Real-Ip")
	if realIP != "" {
		fmt.Println(realIP)
	}

	// Use RemoteAddr if no headers are set
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Println(r.RemoteAddr)
	}

	fmt.Println(ip)
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

		// Add Migada with password Perera and admin as true
		addUserSQL := `INSERT INTO users (username, password, admin) VALUES (?, ?, ?)`
		_, err = db.Exec(addUserSQL, "Migada", "Perera", true)
		if err != nil {
			return nil, fmt.Errorf("failed to add user: %v", err)

		}

		createTableSQL = `CREATE TABLE IF NOT EXISTS sessions (
			session_token TEXT PRIMARY KEY,
			username TEXT NOT NULL
		);`

		_, err = db.Exec(createTableSQL)
		if err != nil {
			return nil, fmt.Errorf("could not create sessions: %v", err)
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

func setImageMeta(db *sql.DB, img ImageMeta) error {
	updateSQL := `UPDATE image_metadata SET description = ?, iso = ?, shutterspeed = ?, aperture = ?, location = ? WHERE filepath = ?`
	_, err := db.Exec(updateSQL, img.Description, img.ISO, img.ShutterSpeed, img.Aperture, img.Location, img.FilePath)
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

func addCookie(db *sql.DB, sessionToken string, username string) error {
	insertSQL := `INSERT INTO sessions (session_token, username) VALUES (?, ?)`
	_, err := db.Exec(insertSQL, sessionToken, username)
	return err
}

func getCookieName(db *sql.DB, sessionToken string) (string, error) {
	querySQL := `SELECT username FROM sessions WHERE session_token = ?`
	row := db.QueryRow(querySQL, sessionToken)
	var username string
	err := row.Scan(&username)
	return username, err
}

func removeCookie(db *sql.DB, sessionToken string) error {
	deleteSQL := `DELETE FROM sessions WHERE session_token = ?`
	_, err := db.Exec(deleteSQL, sessionToken)
	return err
}

func removeCookieByName(db *sql.DB, username string) error {
	deleteSQL := `DELETE FROM sessions WHERE username = ?`
	_, err := db.Exec(deleteSQL, username)
	return err
}

func containsCookie(db *sql.DB, sessionToken string) bool {
	fmt.Println("inside the contains cookie function", sessionToken)
	querySQL := `SELECT session_token FROM sessions WHERE session_token = ?`
	row := db.QueryRow(querySQL, sessionToken)
	var token string
	err := row.Scan(&token)
	if err != nil {
		return false
		fmt.Println("false")
	}
	return true
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

func createUrl(goinTo string, path string) string {
	temp := url.URL{
		Path: goinTo,
	}
	temp.RawQuery = url.Values{
		"data": []string{path},
	}.Encode()

	return temp.String()
}

func descriptionCheck(desc string) bool {
	if len(desc) > 50 {
		return false
	}
	return true
}

func isoCheck(iso string) bool {
	if len(iso) > 6 {
		return false
	}
	isoreg := `^\d{2,6}$`
	re := regexp.MustCompile(isoreg)
	if !re.MatchString(iso) {
		return false
	}
	return true
}

func shutterspeedCheck(ss string) bool {
	if len(ss) > 7 {
		return false
	}
	ssreg := `^(\d+|1/\d{1,5})$`
	re := regexp.MustCompile(ssreg)
	if !re.MatchString(ss) {
		return false
	}
	return true
}

func apertureCheck(apt string) bool {
	if len(apt) > 3 {
		return false
	}
	aptreg := `^\d\.\d+$`
	re := regexp.MustCompile(aptreg)
	if !re.MatchString(apt) {
		return false
	}
	return true
}

func locationCheck(loc string) bool {
	if len(loc) > 340 {
		return false
	}
	return true
}

func applyOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 2:
		return imaging.FlipH(img)
	case 3:
		return imaging.Rotate180(img)
	case 4:
		return imaging.FlipV(img)
	case 5:
		return imaging.Transverse(img)
	case 8:
		return imaging.Rotate90(img)
	default:
		return img
	}
}

func unameCheck(un string) bool {
	return len(un) <= 10
}

func passwordCheck(pw string) bool {
	if len(pw) >= 8 {
		return false
	}
	caps := false
	num := false
	for _, val := range pw {
		if unicode.IsUpper(val) {
			caps = true
		}

		if unicode.IsNumber(val) {
			num = true
		}
	}

	if caps || num {
		return true
	}

	return false
}

func main() {
	db, err := initDB("images.db")
	if err != nil {
		fmt.Println("failed to open database.")
	}

	fs := http.FileServer(http.Dir("images"))
	// http.Handle("/images/", http.StripPrefix("/images/", fs))

	// http.HandleFunc()

	http.Handle("/images/", http.StripPrefix("/images/", fs))

	fs2 := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs2))

	fmt.Println("Starting Server")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		printIP(r)
		data := r.URL.Query().Get("data")
		redir := "/main"
		if data != "" {
			redir = data
		}
		fmt.Println("THe redir" + redir)
		t := template.Must(template.ParseFiles("index.html"))
		t.Execute(w, map[string]string{"data": redir})
	})

	http.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		printIP(r)
		if (r.Header.Get("HX-Request")) != "true" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		tpl := template.Must(template.ParseFiles("main.html"))
		data, err := getAllImageMeta(db)
		if err != nil {
			fmt.Println("Failed to get all data")
		}
		tpl.Execute(w, data)
	})

	http.HandleFunc("/detail", func(w http.ResponseWriter, r *http.Request) {
		printIP(r)
		u := r.URL.String()
		if (r.Header.Get("HX-Request")) != "true" {
			http.Redirect(w, r, createUrl("/", u), http.StatusFound)
			return
		}

		c, err := r.Cookie("session_token")
		if err != nil || !containsCookie(db, c.Value) {
			http.Redirect(w, r, createUrl("/", createUrl("/login", u)), http.StatusFound)
			return
		}

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

	http.HandleFunc("/admin/", func(w http.ResponseWriter, r *http.Request) {
		printIP(r)
		if (r.Header.Get("HX-Request")) != "true" {
			http.Redirect(w, r, createUrl("/", "/admin/"), http.StatusFound)
			return
		}
		fmt.Println("here")
		c, err := r.Cookie("session_token")
		if err != nil || !containsCookie(db, c.Value) {
			http.Redirect(w, r, createUrl("/login", "/admin/"), http.StatusFound)
			return
		}
		uname, err := getCookieName(db, c.Value)
		if err != nil {
			fmt.Println("Failed to get cookie's user" + err.Error())
		}
		fmt.Println(uname)

		admin, err := checkUserAdmin(db, uname)
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
		printIP(r)
		if (r.Header.Get("HX-Request")) != "true" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
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

		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Error reading file", http.StatusMethodNotAllowed)
			return
		}
		defer file.Close()

		buf := new(bytes.Buffer)
		io.Copy(buf, file)

		img, _, err := image.Decode(buf)
		if err != nil {
			http.Error(w, "Error decoding image", http.StatusInternalServerError)
			return
		}

		buf.Reset()
		file.Seek(0, io.SeekStart)
		io.Copy(buf, file)

		ex, err := exif.Decode(buf)
		if err != nil {
			log.Println("No Exif data found, skipping orientation adjustment")
		} else {
			orientationTag, err := ex.Get(exif.Orientation)
			if err == nil {
				orientation, err := orientationTag.Int(0)
				if err == nil {
					img = applyOrientation(img, orientation)
				}
			}
		}

		fmt.Println("done here")

		var webpBuf bytes.Buffer
		err = webp.Encode(&webpBuf, img, &webp.Options{Lossless: false, Quality: 80})
		if err != nil {
			http.Error(w, "Error encoding webp image", http.StatusInternalServerError)
			return
		}
		fmt.Println("done here 2")

		outFilePath := filepath.Join("images", strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))+".webp")
		outFile, err := os.Create(outFilePath)
		if err != nil {
			http.Error(w, "Error saving WebP image", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		fmt.Println("done here 3")

		_, err = outFile.Write(webpBuf.Bytes())
		if err != nil {
			http.Error(w, "Error writing Webp image to file", http.StatusInternalServerError)
			return
		}

		fmt.Println("done here 3")

		fmt.Println(description)
		fmt.Println(iso)
		fmt.Println(shutterSpeed)
		fmt.Println(aperture)
		fmt.Println(outFilePath)
		filepath := outFilePath

		tmpl := template.Must(template.ParseFiles("admin.html"))

		if !descriptionCheck(description) {
			tmpl.ExecuteTemplate(w, "entry-list", ImageMeta{FilePath: filepath, Description: "Failed", ISO: "", ShutterSpeed: "", Aperture: "", Location: "Description invalid: Can't be more than 50 characters long"})
			return
		}

		if !isoCheck(iso) {
			tmpl.ExecuteTemplate(w, "entry-list", ImageMeta{FilePath: filepath, Description: "Failed", ISO: "", ShutterSpeed: "", Aperture: "", Location: "ISO invalid: Example of valid: 400, 25600"})
			return
		}

		if !shutterspeedCheck(shutterSpeed) {
			tmpl.ExecuteTemplate(w, "entry-list", ImageMeta{FilePath: filepath, Description: "Failed", ISO: "", ShutterSpeed: "", Aperture: "", Location: "ShutterSpeed invalid invalid: Example of valid: 1/2000, 300"})
			return
		}

		if !apertureCheck(aperture) {
			tmpl.ExecuteTemplate(w, "entry-list", ImageMeta{FilePath: filepath, Description: "Failed", ISO: "", ShutterSpeed: "", Aperture: "", Location: "Aperture invalid: Example of valid: 4.0, 3.2"})
			return
		}

		if !locationCheck(location) {
			tmpl.ExecuteTemplate(w, "entry-list", ImageMeta{FilePath: filepath, Description: "Failed", ISO: "", ShutterSpeed: "", Aperture: "", Location: "Location invalid: Can't be more than 340 characters long"})
			return
		}

		fmt.Println("done here 6")

		err = addImageMeta(db, ImageMeta{FilePath: filepath, Description: description, ISO: iso, ShutterSpeed: shutterSpeed, Aperture: aperture, Location: location})
		fmt.Println("done", err)
		tmpl.ExecuteTemplate(w, "entry-list", ImageMeta{FilePath: filepath, Description: description, ISO: iso, ShutterSpeed: shutterSpeed, Aperture: aperture, Location: location})
	})

	http.HandleFunc("/admin/update/", func(w http.ResponseWriter, r *http.Request) {
		if (r.Header.Get("HX-Request")) != "true" {
			http.Redirect(w, r, "/admin", http.StatusFound)
			return
		}
		filePath := r.URL.Path[len("/admin/update/"):]
		if r.Method == http.MethodGet {
			fmt.Println("is get in the update")

			data, err := getImageMeta(db, filePath)
			if err != nil {
				fmt.Println("failed to get image's data for initial image")
				return
			}
			t := template.Must(template.ParseFiles("updateform.html"))
			t.Execute(w, data)
		} else if r.Method == http.MethodPost {
			fmt.Println("is post in the update")

			r.ParseForm()
			desc := r.Form.Get("description")
			iso := r.Form.Get("iso")
			ss := r.Form.Get("shutter-speed")
			aperture := r.Form.Get("aperture")
			loc := r.Form.Get("location")

			fmt.Println("description: " + desc)
			fmt.Println("ISO: " + iso)
			fmt.Println("ShutterSpeed: " + ss)
			fmt.Println("Aperture: " + aperture)
			fmt.Println("Location: " + loc)

			if !descriptionCheck(desc) {
				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte("Description Invalid"))
				return
			}

			if !isoCheck(iso) {
				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte("ISO invalid"))
				return
			}

			if !shutterspeedCheck(ss) {
				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte("Shutter Speed invalid"))
				return
			}

			if !apertureCheck(aperture) {
				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte("Shutter Speed invalid"))
				return
			}

			if !locationCheck(loc) {
				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte("Shutter Speed invalid"))
				return
			}

			err := setImageMeta(db, ImageMeta{FilePath: filePath, Description: desc, ISO: iso, ShutterSpeed: ss, Aperture: aperture, Location: loc})
			if err != nil {
				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte("Failed to updata the database entry"))
				return
			}
			http.Redirect(w, r, createUrl("/", "/admin"), http.StatusFound)
		}
	})

	http.HandleFunc("/admin/delete/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		filePath := r.URL.Path[len("/admin/delete/"):]
		fmt.Printf("Deleting item with file path: %s\n", filePath)
		response := "Deleted"
		err := deleteImageMeta(db, filePath)
		if err != nil {
			response = "Failed"
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(response))
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		printIP(r)
		fmt.Println(r.URL.String())
		if r.Method == http.MethodPost {
			fmt.Println("is a post ")
			r.ParseForm()
			uname := r.Form.Get("username")
			pword := r.Form.Get("password")

			t := template.Must(template.ParseFiles("login.html"))

			if !unameCheck(uname) {
				t.ExecuteTemplate(w, "titleb", map[string]string{"data": "Unsucsessful Login"})
				return
			}

			if !passwordCheck(pword) {
				t.ExecuteTemplate(w, "titleb", map[string]string{"data": "Unsucsessful Login"})
				return
			}

			userexists := checkUser(db, uname, pword)

			if !userexists {
				fmt.Println("user doesn't exist")

				t.ExecuteTemplate(w, "titleb", map[string]string{"data": "Unsucsessful Login"})
				return
			}

			err := removeCookieByName(db, uname)

			if err != nil {
				fmt.Println("Failed to remove old cookies")
			}

			cval := convertToSHA256(uname + "@20")

			c := http.Cookie{
				Name:     "session_token",
				Value:    cval,
				HttpOnly: true,
				MaxAge:   7200,
			}

			err = addCookie(db, cval, uname)

			if err != nil {
				fmt.Println("Failed to add cookie to the database")
			}

			http.SetCookie(w, &c)
			fmt.Println("gets to the successfull login")
			path := r.URL.Query().Get("data")
			path, err = url.QueryUnescape(path)
			if err != nil {
				fmt.Println("Failed to get path from url")
			}
			fmt.Println(path)
			// t := template.Must(template.ParseFiles("login.html"))
			// w.Header().Add("HX-Redirect", path)
			// t.ExecuteTemplate(w, "titleb", map[string]string{"data": "Successful Login"})
			http.Redirect(w, r, createUrl("/", path), http.StatusFound)

		} else {
			fmt.Println("is a get")
			if r.Header.Get("HX-Request") != "true" {
				fmt.Println("Is an hx")
				http.Redirect(w, r, createUrl("/", "/login"), http.StatusFound)
			} else {
				fmt.Println("Not an hx")
				path := r.URL.Query().Get("data")
				path, err = url.QueryUnescape(path)
				if err != nil {
					fmt.Println("Failed to get path from url")
				}
				fmt.Println(path)
				t := template.Must(template.ParseFiles("login.html"))
				t.Execute(w, map[string]string{"data": "Login", "redir": path})
			}
		}

	})

	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		printIP(r)
		if (r.Header.Get("HX-Request")) != "true" {
			http.Redirect(w, r, createUrl("/", "/signup"), http.StatusFound)
			return
		}
		if r.Method == http.MethodPost {
			fmt.Println("Inside of signup")
			r.ParseForm()
			uname := r.Form.Get("username")
			pword := r.Form.Get("password")
			tmpl := template.Must(template.ParseFiles("login.html"))

			if !unameCheck(uname) {
				tmpl.ExecuteTemplate(w, "titleb", map[string]string{"data": "Invalid Username"})
				return
			}

			if !passwordCheck(pword) {
				tmpl.ExecuteTemplate(w, "titleb", map[string]string{"data": "Invalid Password"})
				return
			}

			ts := checkUser(db, uname, pword)
			fmt.Println(ts)
			if !ts {
				fmt.Println("User not already exists")
				err := addUser(db, User{Username: uname, Password: pword, admin: false})
				if err != nil {
					fmt.Println("Unsuccessful Sign Up")
					fmt.Println(err)
					ss := err.Error() == "UNIQUE constraint failed: users.username"
					if ss {
						tmpl.ExecuteTemplate(w, "titleb", map[string]string{"data": "Username Already Exists"})
					} else {
						tmpl.ExecuteTemplate(w, "titleb", map[string]string{"data": "Something wrong with the database"})
					}
					return
				}

				tmpl := template.Must(template.ParseFiles("signup.html"))
				w.Header().Set("HX-Redirect", "/login")
				tmpl.ExecuteTemplate(w, "titleb", map[string]string{"data": "Successful Sign Up"})
				return
			}

		} else {
			if (r.Header.Get("HX-Request")) != "true" {
				http.Redirect(w, r, createUrl("/", "/signup"), http.StatusFound)
				return
			}
			t := template.Must(template.ParseFiles("signup.html"))
			t.Execute(w, map[string]string{"data": "Sign Up"})
		}
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		printIP(r)
		c, err := r.Cookie("session_token")
		if err != nil {
			fmt.Println("Failed to get cookie")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		c.MaxAge = -1
		c.Expires = time.Unix(0, 0)

		http.SetCookie(w, c)
		err = removeCookie(db, c.Value)
		if err != nil {
			fmt.Println("Failed to remove cookie")
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	})

	log.Fatal(http.ListenAndServe(":9000", nil))
}
