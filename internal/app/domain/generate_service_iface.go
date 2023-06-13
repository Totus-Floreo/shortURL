package domain

type IGenerateLinkService interface {
	GenerateShortLink() (string, int64)
}
