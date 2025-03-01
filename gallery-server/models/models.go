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
