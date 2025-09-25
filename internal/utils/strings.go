package utils

import "strings"

func FirstChatToLowerCase(str string) string {
	firstChar := str[:1]
	return strings.ToLower(firstChar) + str[1:]
}
