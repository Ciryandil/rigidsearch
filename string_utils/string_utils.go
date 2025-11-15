package string_utils

import "strings"

func CleanWord(word string) string {
	state := 0
	endPos := -1
	startPos := -1
	word = strings.ToLower(word)
	for i := range word {
		if (word[i] >= 'a' && word[i] <= 'z') || (word[i] >= '0' && word[i] <= '9') {
			if state == 0 {
				state = 1
				startPos = i
			} else if state == 2 {
				state = 1
				endPos = -1
			}
		} else {
			if state == 1 {
				state = 2
				endPos = i - 1
			}
		}
	}
	if startPos == -1 {
		return ""
	}
	if endPos == -1 {
		endPos = len(word) - 1
	}
	return word[startPos : endPos+1]
}
