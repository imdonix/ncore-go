package ncore

import "errors"

var (
	ErrNotLoggedIn     = errors.New("not logged in")
	ErrLoginFailed     = errors.New("login failed")
	ErrDownloadFailed  = errors.New("download failed")
	ErrParserError     = errors.New("parser error")
	ErrConnectionError = errors.New("connection error")
)
