package service

import "regexp"

func CheckLink(link string) bool {

	pattern := `^(https?://)?[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(:[0-9]+)?(/.*)?$`

	matched, err := regexp.MatchString(pattern, link)
	if err != nil {
		return false
	}

	return matched
}
