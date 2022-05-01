package utils

import (
	"strings"
	"unicode"
)

func ClearArrString(list []string) []string {
	new := []string{}
	for _, v := range list {
		if len(v) > 0 {
			v = strings.Replace(v, "\t", "", -1)
			new = append(new, strings.Trim(v," "))
		}
	}
	return new
}

func IsEngByLoop(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func GetNameUp (text string) string {
	textUp := ""
	isContinue := false
	if strings.Contains(text, "{") {
		return textUp
	}
	for i, v := range text {
		if isContinue == true {
			isContinue = false
			continue
		}
		if i == 0 {
			textUp = textUp + strings.ToUpper(string(v))
		} else {
			if string(v) == "-" {
				i = i + 1
				textUp = textUp + strings.ToUpper(string(text[i]))
				isContinue = true
			} else {
				textUp = textUp + string(v)
			}
		}
	}
	return textUp
}

func SearchArrString(a []string, b string) bool {
	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}
