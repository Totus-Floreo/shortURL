package service

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	Length           = 10
	LowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	UppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers          = "0123456789"
	Underscore       = "_"
)

type GenerateLinkService struct {
	Mux *sync.Mutex
}

func NewGenerateLinkService() *GenerateLinkService {
	return &GenerateLinkService{
		Mux: new(sync.Mutex),
	}
}

func (s *GenerateLinkService) GenerateShortLink() (string, int64) {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	var builder strings.Builder
	unix := time.Now().Unix()
	chars := []rune(LowercaseLetters + UppercaseLetters + Numbers + Underscore)

	seed := rand.NewSource(unix)
	random := rand.New(seed)

	for i := 0; i < Length; i++ {
		builder.WriteRune(chars[random.Intn(len(chars))])
	}

	return builder.String(), unix
}
