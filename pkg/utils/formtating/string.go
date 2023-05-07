package formtating

import "strings"

// StrPtr strPtr - утилита для создания указателей на строки
func StrPtr(s string) *string {
	return &s
}

// EscapeMarkdownV2 делает текст безопасным для MarkdownV2
func EscapeMarkdownV2(text string) string {
	charsToEscape := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range charsToEscape {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}
	return text
}

// UnescapeMarkdownV2 переделывает текст из MarkdownV2 в обычный
func UnescapeMarkdownV2(text string) string {
	charsToUnescape := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range charsToUnescape {
		text = strings.ReplaceAll(text, "\\"+char, char)
	}
	return text
}
