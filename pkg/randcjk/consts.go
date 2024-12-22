package randcjk

const (
	CJKChinese  = 0b00000001
	CJKJapanese = 0b00000010
	CJKKorean   = 0b00000100
	ASCIILower  = 0b00010000
	ASCIIUpper  = 0b00100000
	ASCIIDigit  = 0b01000000
	ASCIISymbol = 0b10000000

	CJK           = CJKChinese | CJKJapanese | CJKKorean
	ASCIIURL      = ASCIILower | ASCIIDigit
	ASCIIReadable = ASCIILower | ASCIIUpper | ASCIIDigit
	ASCII         = ASCIILower | ASCIIUpper | ASCIIDigit | ASCIISymbol
	CJKASCII      = CJK | ASCII

	CJKChineseStart   = 0x4E00
	CJKChineseEnd     = 0x9FFF
	CJKJapaneseStart  = 0x3040
	CJKJapaneseEnd    = 0x30FF
	CJKKoreanStart    = 0xAC00
	CJKKoreanEnd      = 0xD7A3
	ASCIILowerStart   = 0x61
	ASCIILowerEnd     = 0x7A
	ASCIIUpperStart   = 0x41
	ASCIIUpperEnd     = 0x5A
	ASCIIDigitStart   = 0x30
	ASCIIDigitEnd     = 0x39
	ASCIISymbol1Start = 0x21
	ASCIISymbol1End   = 0x2E // Exclude 0x2F (slash)
	ASCIISymbol2Start = 0x3A
	ASCIISymbol2End   = 0x40
	ASCIISymbol3Start = 0x5B
	ASCIISymbol3End   = 0x60
)
