// Package domain содержит логику запросов (валидации и т.д.) внутри структур Service.
package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
	ErrInternal      = errors.New("internal error")
)

type repo[T any] interface {
	Add(context.Context, T) error
	Delete(context.Context, uuid.UUID) error
	Update(context.Context, T) error
	GetByID(context.Context, uuid.UUID) (T, error)
	FindCopy(context.Context, T) (bool, error)
}

type TxManager interface {
	RunSerializable(context.Context, func(context.Context) error) error
	RunReadUncommited(context.Context, func(context.Context) error) error
}
