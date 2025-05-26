package services

import "errors"

var ErrNotFound = errors.New("resource not found")
var ErrAlreadyExists = errors.New("resource already exists")
