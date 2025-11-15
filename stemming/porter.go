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

func porterStemmerStep2(word string) string {
	if len(word) < 2 {
		return word
	}
	switch word[len(word)-2] {
	case 'a':
		if strings.HasSuffix(word, "ational") && porterPatternCounter(word[:len(word)-7]) > 0 {
			return word[:len(word)-7] + "ate"
		}
		if strings.HasSuffix(word, "tional") && porterPatternCounter(word[:len(word)-6]) > 0 {
			return word[:len(word)-6] + "tion"
		}
	case 'c':
		if strings.HasSuffix(word, "enci") && porterPatternCounter(word[:len(word)-4]) > 0 {
			return word[:len(word)-4] + "ence"
		}
		if strings.HasSuffix(word, "anci") && porterPatternCounter(word[:len(word)-4]) > 0 {
			return word[:len(word)-4] + "ance"
		}
	case 'e':
		if strings.HasSuffix(word, "izer") && porterPatternCounter(word[:len(word)-4]) > 0 {
			return word[:len(word)-4] + "ize"
		}
	case 'l':
		if strings.HasSuffix(word, "abli") && porterPatternCounter(word[:len(word)-4]) > 0 {
			return word[:len(word)-4] + "able"
		}
		if strings.HasSuffix(word, "alli") && porterPatternCounter(word[:len(word)-4]) > 0 {
			return word[:len(word)-4] + "al"
		}
		if strings.HasSuffix(word, "entli") && porterPatternCounter(word[:len(word)-5]) > 0 {
			return word[:len(word)-5] + "ent"
		}
		if strings.HasSuffix(word, "eli") && porterPatternCounter(word[:len(word)-3]) > 0 {
			return word[:len(word)-3] + "el"
		}
		if strings.HasSuffix(word, "ousli") && porterPatternCounter(word[:len(word)-5]) > 0 {
			return word[:len(word)-5] + "ous"
		}
	case 'o':
		if strings.HasSuffix(word, "ization") && porterPatternCounter(word[:len(word)-7]) > 0 {
			return word[:len(word)-7] + "ize"
		}
		if strings.HasSuffix(word, "ation") && porterPatternCounter(word[:len(word)-5]) > 0 {
			return word[:len(word)-5] + "ate"
		}
		if strings.HasSuffix(word, "ator") && porterPatternCounter(word[:len(word)-4]) > 0 {
			return word[:len(word)-4] + "ate"
		}
	case 's':
		if strings.HasSuffix(word, "alism") && porterPatternCounter(word[:len(word)-5]) > 0 {
			return word[:len(word)-5] + "al"
		}
		if strings.HasSuffix(word, "iveness") && porterPatternCounter(word[:len(word)-7]) > 0 {
			return word[:len(word)-7] + "ive"
		}
		if strings.HasSuffix(word, "fulness") && porterPatternCounter(word[:len(word)-7]) > 0 {
			return word[:len(word)-7] + "ful"
		}
		if strings.HasSuffix(word, "ousness") && porterPatternCounter(word[:len(word)-7]) > 0 {
			return word[:len(word)-7] + "ous"
		}
	case 't':
		if strings.HasSuffix(word, "aliti") && porterPatternCounter(word[:len(word)-5]) > 0 {
			return word[:len(word)-5] + "al"
		}
		if strings.HasSuffix(word, "iviti") && porterPatternCounter(word[:len(word)-5]) > 0 {
			return word[:len(word)-5] + "ive"
		}
		if strings.HasSuffix(word, "biliti") && porterPatternCounter(word[:len(word)-6]) > 0 {
			return word[:len(word)-6] + "ble"
		}
	}
	return word
}

func porterStemmerStep3(word string) string {
	if strings.HasSuffix(word, "icate") && porterPatternCounter(word[:len(word)-5]) > 0 {
		return word[:len(word)-5] + "ic"
	}
	if strings.HasSuffix(word, "ative") && porterPatternCounter(word[:len(word)-5]) > 0 {
		return word[:len(word)-5]
	}
	if strings.HasSuffix(word, "alize") && porterPatternCounter(word[:len(word)-5]) > 0 {
		return word[:len(word)-5] + "al"
	}
	if strings.HasSuffix(word, "iciti") && porterPatternCounter(word[:len(word)-5]) > 0 {
		return word[:len(word)-5] + "ic"
	}
	if strings.HasSuffix(word, "ical") && porterPatternCounter(word[:len(word)-4]) > 0 {
		return word[:len(word)-4] + "ic"
	}
	if strings.HasSuffix(word, "ful") && porterPatternCounter(word[:len(word)-3]) > 0 {
		return word[:len(word)-3]
	}
	if strings.HasSuffix(word, "ness") && porterPatternCounter(word[:len(word)-4]) > 0 {
		return word[:len(word)-4]
	}
	return word
}

func porterStemmerStep4(word string) string {
	if strings.HasSuffix(word, "al") && porterPatternCounter(word[:len(word)-2]) > 1 {
		return word[:len(word)-2]
	}
	if strings.HasSuffix(word, "ance") && porterPatternCounter(word[:len(word)-4]) > 1 {
		return word[:len(word)-4]
	}
	if strings.HasSuffix(word, "ence") && porterPatternCounter(word[:len(word)-4]) > 1 {
		return word[:len(word)-4]
	}
	if strings.HasSuffix(word, "er") && porterPatternCounter(word[:len(word)-2]) > 1 {
		return word[:len(word)-2]
	}
	if strings.HasSuffix(word, "ic") && porterPatternCounter(word[:len(word)-2]) > 1 {
		return word[:len(word)-2]
	}
	if strings.HasSuffix(word, "able") && porterPatternCounter(word[:len(word)-4]) > 1 {
		return word[:len(word)-4]
	}
	if strings.HasSuffix(word, "ible") && porterPatternCounter(word[:len(word)-4]) > 1 {
		return word[:len(word)-4]
	}
	if strings.HasSuffix(word, "ant") && porterPatternCounter(word[:len(word)-3]) > 1 {
		return word[:len(word)-3]
	}
	if strings.HasSuffix(word, "ement") && porterPatternCounter(word[:len(word)-5]) > 1 {
		return word[:len(word)-5]
	}
	if strings.HasSuffix(word, "ment") && porterPatternCounter(word[:len(word)-4]) > 1 {
		return word[:len(word)-4]
	}
	if strings.HasSuffix(word, "ent") && porterPatternCounter(word[:len(word)-3]) > 1 {
		return word[:len(word)-3]
	}
	if strings.HasSuffix(word, "ion") && porterPatternCounter(word[:len(word)-3]) > 1 &&
		(strings.HasSuffix(word[:len(word)-3], "t") || strings.HasSuffix(word[:len(word)-3], "s")) {
		return word[:len(word)-2]
	}
	if strings.HasSuffix(word, "ou") && porterPatternCounter(word[:len(word)-2]) > 1 {
		return word[:len(word)-2]
	}
	if strings.HasSuffix(word, "ism") && porterPatternCounter(word[:len(word)-3]) > 1 {
		return word[:len(word)-3]
	}
	if strings.HasSuffix(word, "ate") && porterPatternCounter(word[:len(word)-3]) > 1 {
		return word[:len(word)-3]
	}
	if strings.HasSuffix(word, "iti") && porterPatternCounter(word[:len(word)-3]) > 1 {
		return word[:len(word)-3]
	}
	if strings.HasSuffix(word, "ous") && porterPatternCounter(word[:len(word)-3]) > 1 {
		return word[:len(word)-3]
	}
	if strings.HasSuffix(word, "ive") && porterPatternCounter(word[:len(word)-3]) > 1 {
		return word[:len(word)-3]
	}
	if strings.HasSuffix(word, "ize") && porterPatternCounter(word[:len(word)-3]) > 1 {
		return word[:len(word)-3]
	}
	return word
}

func porterStemmerStep5a(word string) string {
	if strings.HasSuffix(word, "e") {
		patternCount := porterPatternCounter(word[:len(word)-1])
		if patternCount > 1 || (patternCount == 1 && !porterRuleStarO(word[:len(word)-1])) {
			return word[:len(word)-1]
		}
	}
	return word
}

func porterStemmerStep5b(word string) string {
	if len(word) >= 2 && porterPatternCounter(word) > 1 && word[len(word)-1] == 'l' && word[len(word)-2] == 'l' {
		return word[:len(word)-1]
	}
	return word
}

func PorterStemmer(word string) (stem string) {
	word = porterStemmerStep1a(word)
	word = porterStemmerStep1b(word)
	word = porterStemmerStep1c(word)
	word = porterStemmerStep2(word)
	word = porterStemmerStep3(word)
	word = porterStemmerStep4(word)
	word = porterStemmerStep5a(word)
	word = porterStemmerStep5b(word)

	return word
}
