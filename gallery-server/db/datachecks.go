package db

import (
	"fmt"
	"gallery-server/models"
	"regexp"
	"unicode"
)

func DescriptionCheck(desc string) bool {
	if len(desc) > 50 {
		return false
	}
	return true
}

func IsoCheck(iso string) bool {
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

func ShutterspeedCheck(ss string) bool {
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

func ApertureCheck(apt string) bool {
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

func LocationCheck(loc string) bool {
	if len(loc) > 340 {
		return false
	}
	return true
}

func UnameCheck(un string) bool {
	return len(un) <= 10
}

func PasswordCheck(pw string) bool {
	if len(pw) < 8 {
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

func CameraCheck(img *models.ImageMeta) bool {
	if !ApertureCheck(img.Aperture) {
		fmt.Println("aperture fail")
		return false
	}
	if !ShutterspeedCheck(img.ShutterSpeed) {
		fmt.Println("ss fail")
		return false
	}
	if !DescriptionCheck(img.Description) {
		fmt.Println("desc failed")
		return false
	}
	if !IsoCheck(img.ISO) {
		fmt.Println("iso fail")
		return false
	}
	if !LocationCheck(img.Location) {
		return false
	}
	return true
}

func UserCheck(usr *models.User) bool {
	if !UnameCheck(usr.Username) {
		return false
	}
	if !PasswordCheck(usr.Password) {
		return false
	}
	return true
}
