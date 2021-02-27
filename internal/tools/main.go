package tools

import "strings"

func MultiplySpaces(num int) string {
	return MultiplyString(num, " ")
}

func MultiplyString(num int, text string) string {
	var result string
	for i := 0; i <= num; i++ {
		result += text
	}
	return result
}

func GetRateString(text string) string {
	text = text[GetFirstNumberPosition(text):]
	return strings.Replace(text,",",".",1)
}

func GetFirstNumberPosition(text string) int {
	for i:=len(text)-1;i>=0;i-- {
		ascii := int(text[i])
		if ascii >= 48 && ascii <= 57 || ascii==44 {
			continue
		}
		return i+1
	}
	return -1
}