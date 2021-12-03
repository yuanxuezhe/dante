package comm

import (
	"reflect"
	"strconv"
	"unicode/utf8"
)

// Use simple []byte instead of bytes.Buffer to avoid large dependency.
type buffer []byte

func (b *buffer) write(p []byte) {
	*b = append(*b, p...)
}

func (b *buffer) writeString(s string) {
	*b = append(*b, s...)
}

func (b *buffer) writeByte(c byte) {
	*b = append(*b, c)
}

func (bp *buffer) writeRune(r rune) {
	if r < utf8.RuneSelf {
		*bp = append(*bp, byte(r))
		return
	}

	b := *bp
	n := len(b)
	for n+utf8.UTFMax > cap(b) {
		b = append(b, 0)
	}
	w := utf8.EncodeRune(b[n:n+utf8.UTFMax], r)
	*bp = b[:n+w]
}

func doPrintf(format string, a []interface{}) []byte {
	var buff buffer
	end := len(format)

	for i := 0; i < end; {
		lasti := i
		for i < end && format[i] != '%' {
			i++
		}
		if i > lasti {
			buff.writeString(format[lasti:i])
		}
		if i >= end {
			// done processing format string
			break
		}

		i++

		// Process one verb
		k := reflect.TypeOf(a[format[i]-'0'-1]).Kind().String()
		if k == "int" {
			buff.writeString(strconv.Itoa(a[format[i]-'0'-1].(int)))
		} else {
			buff.writeString("'" + a[format[i]-'0'-1].(string) + "'")
		}
		i++
	}

	return buff

}

// Sprintf formats according to a format specifier and returns the resulting string.
func Sprintf(format string, a ...interface{}) string {
	s := doPrintf(format, a)
	return string(s)
}
