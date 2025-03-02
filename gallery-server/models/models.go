package models

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
	Admin    bool
}

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message  string `json:"message"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
}
