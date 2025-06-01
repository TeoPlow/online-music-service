package validation

import (
	"github.com/TeoPlow/online-music-service/src/auth/internal/config"
	"github.com/TeoPlow/online-music-service/src/auth/internal/logger"
	auth "github.com/TeoPlow/online-music-service/src/auth/pkg/authpb"
)

type Validator interface {
	ValidateRegisterUserRequest(req *auth.RegisterUserRequest) error
	ValidateRegisterArtistRequest(req *auth.RegisterArtistRequest) error
	ValidateLoginRequest(req *auth.LoginRequest) error
	ValidateGetUserRequest(req *auth.GetUserRequest) error
	ValidateGetArtistRequest(req *auth.GetArtistRequest) error
	ValidateUpdateUserRequest(req *auth.UpdateUserRequest) error
	ValidateUpdateArtistRequest(req *auth.UpdateArtistRequest) error
	ValidateChangePasswordRequest(req *auth.ChangePasswordRequest) error
	ValidateRefreshTokenRequest(req *auth.RefreshTokenRequest) error
	ValidateLogoutRequest(req *auth.LogoutRequest) error
}

type validator struct {
	countries map[string]bool
}

func NewValidator(cfg *config.Config) Validator {
	countries, err := LoadCountries(cfg.StaticFilesPath)
	if err != nil {
		logger.Log.Error("failed to load countries", "error", err.Error())
		countries = make(map[string]bool)
	}
	return &validator{countries: countries}
}

func (v *validator) ValidateRegisterUserRequest(req *auth.RegisterUserRequest) error {
	return ValidateRegisterUserRequest(req, v.countries)
}

func (v *validator) ValidateRegisterArtistRequest(req *auth.RegisterArtistRequest) error {
	return ValidateRegisterArtistRequest(req, v.countries)
}

func (v *validator) ValidateLoginRequest(req *auth.LoginRequest) error {
	return ValidateLoginRequest(req)
}

func (v *validator) ValidateGetUserRequest(req *auth.GetUserRequest) error {
	return ValidateGetUserRequest(req)
}

func (v *validator) ValidateUpdateUserRequest(req *auth.UpdateUserRequest) error {
	return ValidateUpdateUserRequest(req)
}
func (v *validator) ValidateUpdateArtistRequest(req *auth.UpdateArtistRequest) error {
	return ValidateUpdateArtistRequest(req)
}

func (v *validator) ValidateChangePasswordRequest(req *auth.ChangePasswordRequest) error {
	return ValidateChangePasswordRequest(req)
}

func (v *validator) ValidateRefreshTokenRequest(req *auth.RefreshTokenRequest) error {
	return ValidateRefreshTokenRequest(req)
}

func (v *validator) ValidateLogoutRequest(req *auth.LogoutRequest) error {
	return ValidateLogoutRequest(req)
}

func (v *validator) ValidateGetArtistRequest(req *auth.GetArtistRequest) error {
	return ValidateGetArtistRequest(req)
}
