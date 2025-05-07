// Package domain содержит базовые ошибки и доменные сущности для auth-сервиса.
package domain

import "errors"

var ErrUserNotFound = errors.New("user not found")
