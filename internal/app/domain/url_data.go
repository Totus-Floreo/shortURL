package domain

type URLData struct {
	URLShort string `json:"short"`
	URLLong         // original link struct
	// ExpiresAt time.Time // unrealized, cuz its not specified in test task
}

func NewURLData(short string, long string, addedAt int64) *URLData {
	return &URLData{
		URLShort: short,
		URLLong: URLLong{
			LongURL: long,
			AddedAt: addedAt,
		},
	}
}
