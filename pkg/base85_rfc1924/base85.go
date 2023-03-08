package base85_rfc1924

import (
	"bufio"
	"io"
	"strconv"
)

// 0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!#$%&()*+-;<=>?@^_`{|}~
var b2cTable = [85]byte{
	// 0     1     2     3     4     5     6     7     8     9     A     B     C     D     E     F
	0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, //0
	0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, //1
	0x57, 0x58, 0x59, 0x5A, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C, //2
	0x6D, 0x6E, 0x6F, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7A, 0x21, 0x23, //3
	0x24, 0x25, 0x26, 0x28, 0x29, 0x2A, 0x2B, 0x2D, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F, 0x40, 0x5E, 0x5F, //4
	0x60, 0x7B, 0x7C, 0x7D, 0x7E,
}

var encode [85]byte
var decode [256]byte

func init() {
	copy(encode[:], b2cTable[:])
	for i := 0; i < len(decode); i++ {
		decode[i] = 0xFF
	}
	for i := 0; i < len(encode); i++ {
		decode[encode[i]] = byte(i)
	}
}

// encodeChunk encodes 4 byte-chunk to 5 byte
// if chunk size is less than 4, then it is padded before convertion.
// return written bytes or error
func encodeChunk(dst, src []byte) int {
	if len(src) == 0 {
		return 0
	}

	//read 4 byte as big-endian uint32 into small endian uint32
	var val uint32
	switch len(src) {
	default:
		val |= uint32(src[3])
		fallthrough
	case 3:
		val |= uint32(src[2]) << 8
		fallthrough
	case 2:
		val |= uint32(src[1]) << 16
		fallthrough
	case 1:
		val |= uint32(src[0]) << 24
	}

	buf := [5]byte{0, 0, 0, 0, 0}

	for i := 4; i >= 0; i-- {
		r := val % 85
		val /= 85
		buf[i] = encode[r]
	}

	m := EncodedLen(len(src))
	copy(dst[:], buf[:m])
	return m
}

var decodeBase = [5]uint32{85 * 85 * 85 * 85, 85 * 85 * 85, 85 * 85, 85, 1}

// decodeChunk decodes 5 byte-chunk to 4 byte
// if chunk size is less than 5, then it is padded before conversion.
// return written bytes and error input index
func decodeChunk(dst, src []byte) (int, int) {
	if len(src) == 0 {
		return 0, 0
	}
	var val uint32
	m := DecodedLen(len(src))
	buf := [5]byte{84, 84, 84, 84, 84}
	for i := 0; i < len(src); i++ {
		e := decode[src[i]]
		if e == 0xFF {
			return 0, i + 1
		}
		buf[i] = e
	}

	for i := 0; i < 5; i++ {
		r := buf[i]
		val += uint32(r) * decodeBase[i]
	}
	//small endian uint32 to big endian uint32 in bytes
	switch m {
	default:
		dst[3] = byte(val & 0xff)
		fallthrough
	case 3:
		dst[2] = byte((val >> 8) & 0xff)
		fallthrough
	case 2:
		dst[1] = byte((val >> 16) & 0xff)
		fallthrough
	case 1:
		dst[0] = byte((val >> 24) & 0xff)
	}
	return m, 0
}

// Encode encodes src into dst, return the bytes written
// The dst must have size of EncodedLen(len(src))
func Encode(dst, src []byte) int {
	n := 0
	for len(src) > 0 {
		// Encode 4-byte chunks directly to output slice
		if len(src) >= 4 {
			// Convert 4 bytes to uint32
			val := uint32(src[0])<<24 | uint32(src[1])<<16 | uint32(src[2])<<8 | uint32(src[3])
			for i := 4; i >= 0; i-- {
				r := val % 85
				val /= 85
				dst[i+n] = b2cTable[r]
			}
			n += 5
			src = src[4:]
		} else {
			// Pad and encode remaining 1-3 bytes
			var buf [4]byte
			copy(buf[:], src)
			for i := len(src); i < 4; i++ {
				buf[i] = 0
			}
			encodeChunk(dst[n:], buf[:])
			n += EncodedLen(len(src))
			src = nil
		}
	}
	return n
}

// EncodeToString returns the base85 encoding of src.
func EncodeToString(src []byte) string {
	buf := make([]byte, EncodedLen(len(src)))
	Encode(buf, src)
	return string(buf)
}

// DecodeString returns the bytes represented by the base85 string s.
func DecodeString(src string) ([]byte, error) {
	buf := make([]byte, DecodedLen(len(src)))
	_, err := Decode(buf, []byte(src))
	return buf, err
}

// Decode decodes src into dst, return the bytes written
// The dst must have size of DecodedLen(len(src))
// An CorruptInputError is returned when invalid character is found in src.
func Decode(dst, src []byte) (int, error) {
	f := 0
	t := 0
	for len(src) > 0 {
		if len(src) < 5 {
			w, err := decodeChunk(dst, src)
			if err > 0 {
				return t, CorruptInputError(f + err)
			}
			return t + w, nil
		}

		_, err := decodeChunk(dst[:4], src[:5])
		if err > 0 {
			return t, CorruptInputError(f + err)
		} else {
			t += 4
			f += 5
			src = src[5:]
			dst = dst[4:]
		}
	}
	return t, nil
}

// EncodedLen returns the length in bytes of the base64 encoding of an input
// buffer of length n.
func EncodedLen(n int) int {
	s := n / 4
	r := n % 4
	if r > 0 {
		return s*5 + 5 - (4 - r)
	} else {
		return s * 5
	}
}

// DecodedLen returns the maximum length in bytes of the decoded data
// corresponding to n bytes of base85-encoded data.
func DecodedLen(n int) int {
	s := n / 5
	r := n % 5
	if r > 0 {
		return s*4 + 4 - (5 - r)
	} else {
		return s * 4
	}
}

type encoder struct {
	w       io.Writer
	bufin   [4]byte
	encoded [5]byte
	fill    int
}

func (e *encoder) Write(p []byte) (n int, err error) {
	for len(p) > 0 {
		// Copy input into e.bufin until it's full
		toCopy := len(e.bufin) - e.fill
		if toCopy > len(p) {
			toCopy = len(p)
		}
		copy(e.bufin[e.fill:], p[:toCopy])
		e.fill += toCopy
		p = p[toCopy:]

		// If e.bufin is full, encode it and write to output
		if e.fill == len(e.bufin) {
			encodeChunk(e.encoded[:], e.bufin[:])
			if _, err = e.w.Write(e.encoded[:]); err != nil {
				return
			}
			n += 4
			e.fill = 0
		}
	}
	return n, nil
}

func (e *encoder) Close() error {
	if e.fill > 0 {
		// Pad the input buffer with zeroes and encode
		var buf [4]byte
		copy(buf[:], e.bufin[:e.fill])
		encodeChunk(e.encoded[:], buf[:])
		if _, err := e.w.Write(e.encoded[:EncodedLen(e.fill)]); err != nil {
			return err
		}
	}
	return nil
}

// NewEncoder returns a stream encoder of w.
// All write to the encoder is encoded into base85 and write to w.
// The writer should call Close() to indicate the end of stream
func NewEncoder(w io.Writer) io.WriteCloser {
	return &encoder{
		w: bufio.NewWriterSize(w, 1000),
	}
}

type decoder struct {
	r       io.Reader
	bufin   [1000]byte
	decoded []byte
	fill    int
	err     error
}

// NewDecoder returns a stream decoder of r.
// All read from the reader will read the base85 encoded string from r and decode it.
func NewDecoder(r io.Reader) io.Reader {
	return &decoder{
		r: bufio.NewReaderSize(r, 1000),
	}
}

func (d *decoder) Read(p []byte) (n int, err error) {
	if d.err != nil {
		return 0, d.err
	}

	for len(p) > 0 {
		// try filling the buffer
		var m int
		m, err = d.r.Read(d.bufin[d.fill:])
		d.fill += m
		if err != nil {
			// no further input, decode and copy into p
			d.decoded = make([]byte, DecodedLen(d.fill))
			if d.err == io.EOF {
				k, err := Decode(d.decoded, d.bufin[:d.fill])
				copy(p, d.decoded[:k])
				n += k
				d.fill -= EncodedLen(k)
				if err != nil {
					d.err = err
					return n, err
				}
			} else {
				k, err := Decode(d.decoded, d.bufin[:d.fill-d.fill%5])
				copy(p, d.decoded[:k])
				n += k
				d.fill -= EncodedLen(k)
				if err != nil {
					d.err = err
					return n, err
				}
			}
			d.err = err
			return n, d.err
		}
		//decode d.fill - d.fill % 5 byte of d.bufin
		chunkedMax := d.fill
		d.fill = d.fill % 5
		chunkedMax -= d.fill
		d.decoded = make([]byte, DecodedLen(chunkedMax))
		var k int
		k, err = Decode(d.decoded, d.bufin[:chunkedMax])
		copy(p, d.decoded[:k])
		p = p[k:]
		n += k
		if err != nil {
			d.err = err
			return n, d.err
		}
	}
	return n, d.err
}

type CorruptInputError int64

func (e CorruptInputError) Error() string {
	return "illegal base85 data at input byte " + strconv.FormatInt(int64(e), 10)
}
