package randcjk

import "math/rand"

type UnicodeRange struct {
	Start int
	End   int
}

func (r *UnicodeRange) Random() rune {
	c := rune(RRange(r.Start, r.End))
	if c == '/' || c == '\\' {
		return r.Random()
	}
	return c
}

func (r *UnicodeRange) Length() int {
	return r.End - r.Start
}

func RUnicodeRanges(ranges ...UnicodeRange) rune {
	totalProb := 0
	probs := make([]int, len(ranges))
	for i, r := range ranges {
		probs[i] = r.Length()
		totalProb += probs[i]
	}
	r := RRange(0, totalProb)
	for i, p := range probs {
		if r < p {
			return ranges[i].Random()
		}
		r -= p
	}
	return ranges[len(ranges)-1].Random()
}

func RChar(flag int) rune {
	ranges := []UnicodeRange{}
	if flag&CJKChinese != 0 {
		ranges = append(ranges, UnicodeRange{CJKChineseStart, CJKChineseEnd})
	}
	if flag&CJKJapanese != 0 {
		ranges = append(ranges, UnicodeRange{CJKJapaneseStart, CJKJapaneseEnd})
	}
	if flag&CJKKorean != 0 {
		ranges = append(ranges, UnicodeRange{CJKKoreanStart, CJKKoreanEnd})
	}
	if flag&ASCIILower != 0 {
		ranges = append(ranges, UnicodeRange{ASCIILowerStart, ASCIILowerEnd})
	}
	if flag&ASCIIUpper != 0 {
		ranges = append(ranges, UnicodeRange{ASCIIUpperStart, ASCIIUpperEnd})
	}
	if flag&ASCIIDigit != 0 {
		ranges = append(ranges, UnicodeRange{ASCIIDigitStart, ASCIIDigitEnd})
	}
	if flag&ASCIISymbol != 0 {
		ranges = append(ranges, UnicodeRange{ASCIISymbol1Start, ASCIISymbol1End})
		ranges = append(ranges, UnicodeRange{ASCIISymbol2Start, ASCIISymbol2End})
		ranges = append(ranges, UnicodeRange{ASCIISymbol3Start, ASCIISymbol3End})
	}
	return RUnicodeRanges(ranges...)
}

func RString(length, flag int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = RChar(flag)
	}
	return string(b)
}

func RRange(start, end int) int {
	return start + rand.Intn(end-start)
}
