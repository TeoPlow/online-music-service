// Package storage реализует доступ к хранилищу данных и обработку ошибок уровня хранилища.
package storage

import "errors"

var (
	ErrNotExists = errors.New("not exists")
	ErrSQL       = errors.New("error in SQL query")
)
