package validation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/TeoPlow/online-music-service/src/auth/internal/logger"
)

var (
	countrieOnce   sync.Once
	countriesCache map[string]bool
	countriesErr   error
)

func LoadCountries(staticFilePath string) (map[string]bool, error) {
	countrieOnce.Do(func() {

		filePath := filepath.Join(staticFilePath, "countries.json")
		data, err := os.ReadFile(filePath)
		if err != nil {
			countriesErr = err
			return
		}
		logger.Log.Info("Loading countries from:", "path", filePath)

		var countries []string
		if err := json.Unmarshal(data, &countries); err != nil {
			countriesErr = err
			return
		}
		countriesCache = make(map[string]bool)
		for _, countryName := range countries {
			countriesCache[strings.ToLower(countryName)] = true
		}
	})
	return countriesCache, countriesErr
}

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return ErrInvalidUsername
	}
	return nil
}

func ValidateAuthor(author string) error {
	if len(author) < 3 || len(author) > 20 {
		return ErrInvalidAuthor
	}
	return nil
}

func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return ErrPasswordMissingUppercase
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return ErrPasswordMissingLowercase
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return ErrPasswordMissingDigit
	}
	if !regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password) {
		return ErrPasswordMissingSpecial
	}
	return nil
}

func ValidateAge(age int32) error {
	if age <= 0 || age > 150 {
		return ErrInvalidAge
	}
	return nil
}

func ValidateRole(role string) error {
	validRoles := map[string]bool{
		"admin":  true,
		"user":   true,
		"artist": true,
	}
	if !validRoles[strings.ToLower(role)] {
		return ErrInvalidRole
	}
	return nil
}

func ValidateCountry(country string, countries map[string]bool) error {
	if country == "" {
		return ErrCountryEmpty
	}
	if countries == nil {
		return ErrInvalidCountry
	}
	if !countries[strings.ToLower(country)] {
		return ErrInvalidCountry
	}
	return nil
}
