package convert

import "regexp"

type Binary struct {
	zero  byte
	one   byte
	lsb   byte
	rsb   byte
	space byte
}

func NewBinary(zero, one, lsb, rsb, space byte) *Binary {
	bin := &Binary{
		zero:  zero,
		one:   one,
		lsb:   lsb,
		rsb:   rsb,
		space: space,
	}
	return bin
}

func Default() *Binary {
	return NewBinary(zero, one, lsb, rsb, space)
}

// ByteToBinaryString get the string in binary format of a byte or uint8.
func (bin *Binary) ByteToBinaryString(b byte) string {
	buf := make([]byte, 0, 8)
	buf = bin.appendBinaryString(buf, b)
	return string(buf)
}

// BytesToBinaryString get the string in binary format of a []byte or []int8.
func (bin *Binary) BytesToBinaryString(bs []byte) string {
	l := len(bs)
	bl := l*8 + l + 1
	buf := make([]byte, 0, bl)
	buf = append(buf, bin.lsb)
	for _, b := range bs {
		buf = bin.appendBinaryString(buf, b)
		buf = append(buf, bin.space)
	}
	buf[bl-1] = bin.rsb
	return string(buf)
}

// regex for delete useless string which is going to be in binary format.
var rbDel = regexp.MustCompile(`[^01]`)

// BinaryStringToBytes get the binary bytes according to the
// input string which is in binary format.
func (bin *Binary) BinaryStringToBytes(s string) (bs []byte) {
	if len(s) == 0 {
		panic(ErrEmptyString)
	}

	s = rbDel.ReplaceAllString(s, "")
	l := len(s)
	if l == 0 {
		panic(ErrBadStringFormat)
	}

	mo := l % 8
	l /= 8
	if mo != 0 {
		l++
	}
	bs = make([]byte, 0, l)
	mo = 8 - mo
	var n uint8
	for i, b := range []byte(s) {
		m := (i + mo) % 8
		switch b {
		case bin.one:
			n += uint8arr[m]
		}
		if m == 7 {
			bs = append(bs, n)
			n = 0
		}
	}
	return
}

// append bytes of string in binary format.
func (bin *Binary) appendBinaryString(bs []byte, b byte) []byte {
	var a byte
	for i := 0; i < 8; i++ {
		a = b
		b <<= 1
		b >>= 1
		switch a {
		case b:
			bs = append(bs, bin.zero)
		default:
			bs = append(bs, bin.one)
		}
		b <<= 1
	}
	return bs
}
