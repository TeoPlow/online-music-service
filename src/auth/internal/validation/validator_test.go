package validation

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/TeoPlow/online-music-service/src/auth/internal/config"
	"github.com/TeoPlow/online-music-service/src/auth/internal/logger"
	auth "github.com/TeoPlow/online-music-service/src/auth/pkg/authpb"
)

func TestValidateRegisterUserRequest(t *testing.T) {
	tempDir := t.TempDir()
	if logger.Log == nil {
		logger.Log = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	testCountries := []string{"russia", "united states of america", "germany", "france"}
	jsonData, err := json.Marshal(testCountries)
	if err != nil {
		t.Fatalf("Failed to marshal test countries: %v", err)
	}

	staticDir := filepath.Join(tempDir, "static")
	if err := os.MkdirAll(staticDir, 0755); err != nil {
		t.Fatalf("Failed to create static directory: %v", err)
	}

	if err := os.WriteFile(filepath.Join(staticDir, "countries.json"), jsonData, 0644); err != nil {
		t.Fatalf("Failed to write test countries file: %v", err)
	}

	countrieOnce = sync.Once{}
	countriesCache = nil
	countriesErr = nil

	cfg := &config.Config{
		StaticFilesPath: staticDir,
	}
	countries, err := LoadCountries(cfg.StaticFilesPath)
	if err != nil {
		t.Fatalf("Failed to load countries: %v", err)
	}
	if len(countries) == 0 {
		t.Fatal("No countries loaded, expected at least one country")
	}

	var validCountry string
	for country := range countries {
		validCountry = country
		break
	}
	tests := []struct {
		name      string
		req       *auth.RegisterUserRequest
		countries map[string]bool
		wantErr   bool
	}{
		{
			name: "valid request",
			req: &auth.RegisterUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "Password123!",
				Gender:   true,
				Country:  validCountry,
				Age:      25,
			},
			countries: countries,
			wantErr:   false,
		},
		{
			name: "invalid username",
			req: &auth.RegisterUserRequest{
				Username: "ab",
				Email:    "test@example.com",
				Password: "Password123!",
				Gender:   true,
				Country:  "RU",
				Age:      25,
			},
			countries: countries,

			wantErr: true,
		},
		{
			name: "invalid email",
			req: &auth.RegisterUserRequest{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "Password123!",
				Gender:   true,
				Country:  "RU",
				Age:      25,
			},
			countries: countries,

			wantErr: true,
		},
		{
			name: "invalid password",
			req: &auth.RegisterUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "short",
				Gender:   true,
				Country:  "RU",
				Age:      25,
			},
			countries: countries,

			wantErr: true,
		},
		{
			name: "invalid age",
			req: &auth.RegisterUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "Password123!",
				Gender:   true,
				Country:  "RU",
				Age:      -1,
			},
			countries: countries,

			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateRegisterUserRequest(test.req, test.countries)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateRegisterRequest() = %v, want %v", err, test.wantErr)
			}
		})
	}
}

func TestValidateLoginRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *auth.LoginRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: &auth.LoginRequest{
				Username: "Stepashka331",
				Password: "Password123!",
			},
			wantErr: false,
		},
		{
			name: "invalid username",
			req: &auth.LoginRequest{
				Username: "_",
				Password: "Password123!",
			},
			wantErr: true,
		},
		{
			name: "invalid password",
			req: &auth.LoginRequest{
				Username: "Stepashka331",
				Password: "short",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateLoginRequest(test.req)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateLoginRequest() = %v, want %v", err, test.wantErr)
			}
		})
	}
}

func TestValidateGetUserRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *auth.GetUserRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: &auth.GetUserRequest{
				Id: "user-id",
			},
			wantErr: false,
		},
		{
			name: "empty id",
			req: &auth.GetUserRequest{
				Id: "",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateGetUserRequest(test.req)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateGetUserRequest() = %v, want %v", err, test.wantErr)
			}
		})
	}
}

func TestValidateUpdateUserRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *auth.UpdateUserRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: &auth.UpdateUserRequest{
				Id:       "user-id",
				Username: "new-username",
				Age:      25,
			},
			wantErr: false,
		},
		{
			name: "empty id",
			req: &auth.UpdateUserRequest{
				Id:       "",
				Username: "new-username",
				Age:      25,
			},
			wantErr: true,
		},
		{
			name: "invalid username",
			req: &auth.UpdateUserRequest{
				Id:       "user-id",
				Username: "ab",
				Age:      25,
			},
			wantErr: true,
		},
		{
			name: "invalid age",
			req: &auth.UpdateUserRequest{
				Id:       "user-id",
				Username: "new-username",
				Age:      -1,
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateUpdateUserRequest(test.req)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateUpdateUserRequest() = %v, want %v", err, test.wantErr)
			}
		})
	}
}

func TestValidateChangePasswordRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *auth.ChangePasswordRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: &auth.ChangePasswordRequest{
				Id:          "user-id",
				OldPassword: "OldPassword123!",
				NewPassword: "NewPassword123!",
			},
			wantErr: false,
		},
		{
			name: "empty id",
			req: &auth.ChangePasswordRequest{
				Id:          "",
				OldPassword: "OldPassword123!",
				NewPassword: "NewPassword123!",
			},
			wantErr: true,
		},
		{
			name: "invalid old password",
			req: &auth.ChangePasswordRequest{
				Id:          "user-id",
				OldPassword: "short",
				NewPassword: "NewPassword123!",
			},
			wantErr: true,
		},
		{
			name: "invalid new password",
			req: &auth.ChangePasswordRequest{
				Id:          "user-id",
				OldPassword: "OldPassword123!",
				NewPassword: "short",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateChangePasswordRequest(test.req)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateChangePasswordRequest() = %v, want %v", err, test.wantErr)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			name:     "valid username",
			username: "testuser",
			wantErr:  false,
		},
		{
			name:     "invalid username",
			username: "ab",
			wantErr:  true,
		},
		{
			name:     "empty username",
			username: "",
			wantErr:  true,
		},
		{
			name:     "username too long",
			username: "thisusernameiswaytoolong",
			wantErr:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateUsername(test.username)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateUsername(%q) = %v, want %v", test.username, err, test.wantErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "valid email",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "invalid email",
			email:   "invalid-email",
			wantErr: true,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
		},
		{
			name:    "missing @",
			email:   "test.com",
			wantErr: true,
		},
		{
			name:    "missing models",
			email:   "test@",
			wantErr: true,
		},
		{
			name:    "missing local part",
			email:   "@example.com",
			wantErr: true,
		},
		{
			name:    "missing @ and models",
			email:   "test",
			wantErr: true,
		},
		{
			name:    "missing @ and local part",
			email:   "test.com",
			wantErr: true,
		},
		{
			name:    "missing @ and local part and models",
			email:   "test",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateEmail(test.email)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateEmail(%q) = %v, want %v", test.email, err, test.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "Password123!",
			wantErr:  false,
		},
		{
			name:     "short password",
			password: "short",
			wantErr:  true,
		},
		{
			name:     "missing uppercase",
			password: "password123!",
			wantErr:  true,
		},
		{
			name:     "missing lowercase",
			password: "PASSWORD123!",
			wantErr:  true,
		},
		{
			name:     "missing number",
			password: "Password!",
			wantErr:  true,
		},
		{
			name:     "missing special character",
			password: "Password123",
			wantErr:  true,
		},
		{
			name:     "missing uppercase, lowercase, number, and special character",
			password: "password",
			wantErr:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidatePassword(test.password)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidatePassword(%q) = %v, want %v", test.password, err, test.wantErr)
			}
		})
	}
}

func TestValidateAge(t *testing.T) {
	tests := []struct {
		name    string
		age     int32
		wantErr bool
	}{
		{
			name:    "valid age",
			age:     25,
			wantErr: false,
		},
		{
			name:    "invalid age",
			age:     151,
			wantErr: true,
		},
		{
			name:    "negative age",
			age:     -1,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateAge(test.age)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateAge(%d) = %v, want %v", test.age, err, test.wantErr)
			}
		})
	}
}

func TestValidateRole(t *testing.T) {
	tests := []struct {
		name    string
		role    string
		wantErr bool
	}{
		{
			name:    "valid role",
			role:    "admin",
			wantErr: false,
		},
		{
			name:    "invalid role",
			role:    "invalid-role",
			wantErr: true,
		},
		{
			name:    "empty role",
			role:    "",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateRole(test.role)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateRole(%q) = %v, want %v", test.role, err, test.wantErr)
			}
		})
	}
}

func TestValidateRegistrationRole(t *testing.T) {
	tests := []struct {
		name    string
		role    string
		wantErr bool
	}{
		{
			name:    "valid user role",
			role:    "user",
			wantErr: false,
		},
		{
			name:    "valid artist role",
			role:    "artist",
			wantErr: false,
		},
		{
			name:    "invalid admin role",
			role:    "admin",
			wantErr: true,
		},
		{
			name:    "invalid role",
			role:    "invalid-role",
			wantErr: true,
		},
		{
			name:    "empty role",
			role:    "",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateRegistrationRole(test.role)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateRegistrationRole(%q) = %v, want %v", test.role, err, test.wantErr)
			}
		})
	}
}
func TestValidateCountry(t *testing.T) {
	countries := map[string]bool{
		"russia":                   true,
		"united states of america": true,
		"germany":                  true,
		"france":                   true,
	}

	tests := []struct {
		name    string
		country string
		wantErr bool
	}{
		{
			name:    "valid country",
			country: "russia",
			wantErr: false,
		},
		{
			name:    "valid country mixed case",
			country: "Russia",
			wantErr: false,
		},
		{
			name:    "valid country uppercase",
			country: "RUSSIA",
			wantErr: false,
		},
		{
			name:    "valid country with spaces",
			country: "united states of america",
			wantErr: false,
		},
		{
			name:    "invalid country",
			country: "not-a-country",
			wantErr: true,
		},
		{
			name:    "empty country",
			country: "",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateCountry(test.country, countries)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateCountry(%q) = %v, want error: %v", test.country, err, test.wantErr)
			}
		})
	}

	t.Run("nil countries map", func(t *testing.T) {
		err := ValidateCountry("russia", nil)
		if err == nil {
			t.Errorf("Expected error with nil countries map, got nil")
		}
	})
}

func TestLoadCountries(t *testing.T) {
	if logger.Log == nil {
		logger.Log = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	countrieOnce = sync.Once{}
	countriesCache = nil
	countriesErr = nil
	cfg := &config.Config{
		StaticFilesPath: t.TempDir(),
	}
	countriesJSON := `["russia", "united states of america", "germany", "france"]`
	if err := os.WriteFile(filepath.Join(cfg.StaticFilesPath, "countries.json"),
		[]byte(countriesJSON), 0644); err != nil {
		t.Fatalf("Failed to write countries.json: %v", err)
	}
	countries, err := LoadCountries(cfg.StaticFilesPath)
	if err != nil {
		t.Fatalf("LoadCountries() error = %v", err)
	}
	if len(countries) == 0 {
		t.Fatal("LoadCountries() returned empty map, expected at least one country")
	}

	expectedCountries := []string{"russia", "united states of america", "germany", "france"}
	for _, country := range expectedCountries {
		foundCountry := false
		for actualCountry := range countries {
			if strings.EqualFold(actualCountry, country) {
				foundCountry = true
				break
			}
		}
		if !foundCountry {
			t.Logf("Note: Expected country %q not found in loaded countries.", country)
		}
	}

	var someCountry string
	for country := range countries {
		someCountry = country
		break
	}

	if someCountry != "" {
		upperCountry := strings.ToUpper(someCountry)
		if err := ValidateCountry(upperCountry, countries); err != nil {
			t.Errorf("ValidateCountry() failed for uppercase variant of %q: %v", someCountry, err)
		}
	}
}
