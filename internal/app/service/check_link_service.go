package service

import "regexp"

const pattern = `^(https?://)?[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(:[0-9]+)?(/.*)?$`

func CheckLink(link string) bool {

	matched, err := regexp.MatchString(pattern, link)
	if err != nil {
		return false
	}

	return matched
}
