package domain

import "errors"

var (
	ErrorInvalidShort    = errors.New("invalid short link")
	ErrorInvalidLink     = errors.New("invalid link")
	ErrorInvalidDecode   = errors.New("link cant decode")
	ErrorLinkNotFound    = errors.New("link not found")
	ErrorGenerateTimeout = errors.New("generate short link timeout")
)
