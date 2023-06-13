package domain

type URLLong struct {
	LongURL string `json:"link"`
	AddedAt int64  `json:"added"` // addition to justify the construction
}
