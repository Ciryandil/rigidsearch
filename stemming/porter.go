package stemming

import "strings"

func isVowel(letter rune, prevLetterPtr *rune) bool {
	if letter == 'a' || letter == 'e' || letter == 'i' || letter == 'o' || letter == 'u' {
		return true
	}
	if letter == 'y' {
		if prevLetterPtr == nil {
			return false
		}
		prevLetter := *prevLetterPtr
		if prevLetter == 'a' || prevLetter == 'e' || prevLetter == 'i' || prevLetter == 'o' || prevLetter == 'u' {
			return false
		} else {
			return true
		}
	}
	return true
}

func containsVowel(word string) bool {
	runeArr := []rune(word)
	var prevLetterPtr *rune = nil
	for itr := range runeArr {
		if itr > 0 {
			prevLetterPtr = &runeArr[itr-1]
		}
		if isVowel(runeArr[itr], prevLetterPtr) {
			return true
		}
	}
	return false
}

func endsInDoubleConsonant(word string) bool {
	runesArr := []rune(word)
	if len(runesArr) >= 2 &&
		!isVowel(runesArr[len(runesArr)-1], &runesArr[len(runesArr)-2]) {
		var prevLetterPtr *rune
		if len(runesArr) >= 3 {
			prevLetterPtr = &runesArr[len(runesArr)-3]
		}
		if !isVowel(runesArr[len(runesArr)-2], prevLetterPtr) {
			return true
		}
	}
	return false
}

func porterRuleStarO(word string) bool {
	runesArr := []rune(word)
	if len(runesArr) >= 3 &&
		runesArr[len(runesArr)-1] != 'w' &&
		runesArr[len(runesArr)-1] != 'x' &&
		runesArr[len(runesArr)-1] != 'y' &&
		!isVowel(runesArr[len(runesArr)-1], &runesArr[len(runesArr)-2]) &&
		isVowel(runesArr[len(runesArr)-2], &runesArr[len(runesArr)-3]) {
		var prevLetterPtr *rune
		if len(runesArr) >= 4 {
			prevLetterPtr = &runesArr[len(runesArr)-4]
		}
		if !isVowel(runesArr[len(runesArr)-3], prevLetterPtr) {
			return true
		}
	}
	return false
}

func porterPatternCounter(word string) int {
	count := 0
	state := 0
	runeArr := []rune(word)
	var prevLetterPtr *rune = nil
	for itr := range runeArr {
		if itr > 0 {
			prevLetterPtr = &runeArr[itr-1]
		}
		if isVowel(runeArr[itr], prevLetterPtr) {
			if state == 0 {
				state = 1
			} else if state == 2 {
				count += 1
				state = 1
			}
		} else {
			if state == 1 {
				state = 2
			}
		}
	}
	if state == 2 {
		count += 1
	}

	return count
}

func porterStemmerStep1a(word string) string {
	if strings.HasSuffix(word, "sses") {
		return word[:len(word)-4] + "ss"
	}
	if strings.HasSuffix(word, "ies") {
		return word[:len(word)-3] + "i"
	}
	if strings.HasSuffix(word, "s") {
		return word[:len(word)-1]
	}
	return word
}

func porterStemmerStep1b(word string) string {
	convertedWord := word
	requiresFurtherProcessing := false
	if strings.HasSuffix(word, "eed") {
		patternCount := porterPatternCounter(word[:len(word)-3])
		if patternCount > 0 {
			return word[:len(word)-3] + "ee"
		}
	}
	if strings.HasSuffix(word, "ed") {
		if containsVowel(word[:len(word)-2]) {
			convertedWord = word[:len(word)-2]
			requiresFurtherProcessing = true
		}
	}
	if strings.HasSuffix(word, "ing") {
		if containsVowel(word[:len(word)-3]) {
			convertedWord = word[:len(word)-3]
			requiresFurtherProcessing = true
		}
	}
	if requiresFurtherProcessing {
		if strings.HasSuffix(convertedWord, "at") || strings.HasSuffix(convertedWord, "bl") {
			convertedWord += "e"
		}
		if endsInDoubleConsonant(convertedWord) && convertedWord[len(convertedWord)-1] != 'l' && convertedWord[len(convertedWord)-1] != 's' && convertedWord[len(convertedWord)-1] != 'z' {
			convertedWord = convertedWord[:len(word)-1]
		}
		if porterRuleStarO(convertedWord) && porterPatternCounter(convertedWord) == 1 {
			convertedWord += "e"
		}
	}
	return convertedWord
}

func porterStemmerStep1c(word string) string {
	if strings.HasSuffix(word, "y") && containsVowel(word[:len(word)-1]) {
		return word[:len(word)-1] + "i"
	}
	return word
}

func PorterStemmer(word string) (stem string) {
	word = porterStemmerStep1a(word)
	word = porterStemmerStep1b(word)
	word = porterStemmerStep1c(word)
	// further to do
	return word
}
