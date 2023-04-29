package formtating

import "strings"

// StrPtr strPtr - утилита для создания указателей на строки
func StrPtr(s string) *string {
	return &s
}

func EscapeMarkdownV2(text string) string {
	charsToEscape := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range charsToEscape {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}
	return text
}
