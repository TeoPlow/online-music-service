package auth

import (
	"github.com/TeoPlow/online-music-service/src/auth/internal/models"
	authproto "github.com/TeoPlow/online-music-service/src/auth/pkg/authpb"
)

func toProtoRole(r models.Role) authproto.Role {
	switch r {
	case models.RoleUser:
		return authproto.Role_USER
	case models.RoleArtist:
		return authproto.Role_ARTIST
	case models.RoleAdmin:
		return authproto.Role_ADMIN
	default:
		return authproto.Role_ROLE_UNSPECIFIED
	}
}
