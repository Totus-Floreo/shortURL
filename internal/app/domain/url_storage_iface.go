package domain

import "context"

type IUrlStorage interface {
	AddUrl(context.Context, URLData) error
	GetUrl(context.Context, string) (*URLLong, error)
}
