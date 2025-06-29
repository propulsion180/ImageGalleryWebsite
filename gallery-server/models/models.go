package models

type ImageMeta struct {
	FilePath     string `json:"filepath"`
	Description  string `json:"description"`
	ISO          string `json:"iso"`
	ShutterSpeed string `json:"shutterspeed"`
	Aperture     string `json:"aperture"`
	Location     string `json:"location"`
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

type SingleImageData struct {
	FilePath string `json:"filepath"`
}
