package utils

import "strings"

func SlugGenerator(text string) string {
	splittedText := strings.Split(text, " ")

	var slug = ""

	for index, text := range splittedText {
		slug += strings.ToLower(text)

		if index != len(splittedText)-1 {
			slug += "-"
		}
	}

	return slug
}
