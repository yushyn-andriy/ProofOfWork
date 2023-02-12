package utilities

import (
	crrand "crypto/rand"
	"log"
	"math/rand"

	"unicode"
	"unsafe"
)

var letterRunes = []rune("!#$%&()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~ĵĶķĸĹĺĻļĽľĿŀŁłŃńŅņŇňŉŊŋŌōŎŏŐőŒœŔŕŖśŜŝŞşŠšŢţŤťŦŧŨũŪūŬŭŮůŰűŲųŴŵŶŷŸŹźŻżŽžΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩαβγδεζηθικλμνξοπρστυφχψωЀЁЂЃЄЅІЇЈЉЊЋЌЍЎЏАБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдежзийклмнопрстуфхцчшщъыьэюяĀāĂăĄąĆćĈĉĊċČčĎďĐđĒēĔĕĖėĘęĚěĜĝĞğĠġĢģĤĥĦħĨĩĪīĬĭĮįİıĲĳĴœŔŕŖŗŘřŚśŜŝŞşŠšŢţŤťŦŧŨũŪūŬŭŮůŰűŲųŴŵŶŷŸŹźŻżŽžſƀƁƂƃƄƅƆƇƈƉƊƋƌƍƎƏƐƑƒƓƔƕƖƗƘƙƚƛƜƝƞƟƠơƢƣƤƥƦƧƨƩƪƫƬƭƮƯưƱƲƳƴƵƶƷƸƹƺƻƼƽƾƿǀǁǂǃǄǅǆǇǈǉǊǋǌǍǎǏǐǑǒǓǔǕǖǗǘǙǚǛǜǝǞǟǠǡǢǣǤǥǦǧǨǩǪǫǬǭǮǯǰǱǲǳǴǵǶǷǸǹǺǻǼǽǾǿȀȁȂȃȄȅȆȇȈȉȊȋȌȍȎ")

const letterBytes = "!#$%&()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~"
const (
	// So for example if we have 52 letters,
	// it requires 6 bits to represent it: 52 = 110100b.
	// 6 bits to represent a letter index
	letterIdxBits = 6

	// All 1-bits, as many as letterIdxBits 1
	// 00000001 -> 00(111111)
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits // # of letter indices fitting in 63 bits
)

// RandStringBytesMaskImprSrcUnsafe generates random string of length n
// from the set of predefined bytes
func RandStringBytesMaskImprSrcUnsafe(n int, src rand.Source) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}

		// In order to understand whats going on
		// commented with some example
		// 00000000000000000111111 letterIdxMask
		// 11101010010001111001000 cache
		// 00000000000000000001000 cache & letterIdxMask
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}

		// 11101010010001111001000 cache
		// 00000011101010010001111 cache >> 6
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

// RandStringRunes generates random string of length n
// from the set of predefined runes
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// RandUnicodeString generates random string of length n
func RandUnicodeString(n int) string {
	characterSet := make([]rune, n)
	i := 0
	for i < n {
		r := rand.Intn(0x10FFFF)
		if unicode.IsPrint(rune(r)) && !unicode.IsSpace(rune(r)) {
			characterSet[i] = rune(r)
			i++
		}
	}
	return string(characterSet)
}

// RandBytes generates random slice of bytes of length n
func RandBytes(n int) []byte {
	b := make([]byte, n)
	_, err := crrand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
