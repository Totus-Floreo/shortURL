package domain

import "context"

type IUrlService interface {
	CreateUrl(context.Context, string) (string, error)
	GetUrl(context.Context, string) (string, error)
}
