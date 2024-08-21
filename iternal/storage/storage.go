package storage

import "errors"

var (
	ErrorNotFound = errors.New("url not found")
	ErrorExists   = errors.New("url already exists")
	ErrorNotExist = errors.New("alias not exist")
)
