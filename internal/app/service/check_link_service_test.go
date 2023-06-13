package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestCaseCheckLink struct {
	Title   string
	Link    string
	Correct bool
}

func TestCheckLink(t *testing.T) {
	var tests = []TestCaseCheckLink{
		TestCaseCheckLink{
			Title:   "Correct simple link w/out http",
			Link:    "example.com",
			Correct: true,
		},
		TestCaseCheckLink{
			Title:   "Correct simple link w/ http",
			Link:    "http://example.com",
			Correct: true,
		},
		TestCaseCheckLink{
			Title:   "Correct link w/ route /123123/c/123",
			Link:    "http://example.com/123123/c/123",
			Correct: true,
		},
		TestCaseCheckLink{
			Title:   "Correct link w/ params /watch?v=dQw4w9WgXcQ",
			Link:    "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			Correct: true,
		},
		TestCaseCheckLink{
			Title:   "Correct link w/ params and overload /results?search_query=super+sonic",
			Link:    "https://www.youtube.com/results?search_query=super+sonic",
			Correct: true,
		},
		TestCaseCheckLink{
			Title:   "Invalid link: missing host",
			Link:    "http://",
			Correct: false,
		},
		TestCaseCheckLink{
			Title:   "Invalid link: unsupported scheme",
			Link:    "ftp://example.com",
			Correct: false,
		},
		TestCaseCheckLink{
			Title:   "Invalid link: query parameter without value",
			Link:    "http://example.com?param=",
			Correct: false,
		},
		TestCaseCheckLink{
			Title:   "Invalid link: missing scheme and host",
			Link:    "",
			Correct: false,
		},
		TestCaseCheckLink{
			Title:   "Invalid link: contains whitespace",
			Link:    "http://example.com /page",
			Correct: false,
		},
		TestCaseCheckLink{
			Title:   "Invalid link: contains special characters",
			Link:    "http://example.com/\npage",
			Correct: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Title, func(t *testing.T) {
			answer := CheckLink(test.Link)

			require.Equal(t, test.Correct, answer)
		})
	}
}
