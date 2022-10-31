package goio

import (
	"errors"
	"io"
	"unicode/utf8"
)

// StringReader是对strings.Reader的完整拷贝，主要是为了注释和学习。
// A Reader 通过读取一个底层的字符串实现了 io.Reader, io.ReaderAt, io.ByteReader, io.ByteScanner,
// io.RuneReader, io.RuneScanner, io.Seeker 和 io.WriterTo 接口。

// Reader的”零值“相当于读取空字符串的读取器。
type StringReader struct {
	s        string
	i        int64 // 表示当前所读取的字节（byte）的序号（未读）
	prevRune int   // 前一个字符（rune）所在的字节序号（小于i），用于字符读取回退操作；或者 < 0，表示上一个操作不是ReadRune。
}

// returns the number of bytes of the unread portion of the string.
// Len 返回字符串中还没有读取部分的长度。
func (r *StringReader) Len() int {
	if r.i >= int64(len(r.s)) {
		return 0
	}
	return int(int64(len(r.s)) - r.i)
}

// Size returns the original length of the underlying string.
// Size is the number of bytes available for reading via ReadAt.
// The returned value is always the same and is not affected by calls
// to any other method.
// Size返回底层字符串初始长度。Size是通过ReadAt可以读取的字节数。
// 无论调用了其他什么方法，Size()方法返回值总是一样 ，不受影响。
func (r *StringReader) Size() int64 { return int64(len(r.s)) }

// Read implements the io.Reader interface.
// Read 方法实现了io.Reader接口
// !!! 该方法每次调用，当前所读取的字节（byte）的序号都会相应地前进。
func (r *StringReader) Read(b []byte) (n int, err error) {
	//!!! 如果当前所读取的字节序号已经超过了最后一个字节的序号，意味着已经读取完毕，返回io.EOF错误。
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	//!!!前一个字符序号序号置为-1，意味着不是按字符读取。
	r.prevRune = -1
	//copy试图将当前序号到后续所有的字节（r.s[r.i:]）都拷贝到b中，
	//!!! copy不会重新分配内存，因此，不会改变目标切片的长度，
	//!!! 最终能拷贝多少，取决于目标切片b的长度与可拷贝字符的长度的最小值。
	n = copy(b, r.s[r.i:])
	r.i += int64(n)
	return
}

// ReadAt implements the io.ReaderAt interface.
// ReadAt实现了io.ReaderAt接口
// !!!  ReadAt 从指定的偏移处读取字节到目标的切片中。
// !!!  这种读取方式不会改变正常读取方式下的当前所读取的字节序号。
func (r *StringReader) ReadAt(b []byte, off int64) (n int, err error) {
	// cannot modify state - see io.ReaderAt
	if off < 0 {
		return 0, errors.New("strings.Reader.ReadAt: negative offset")
	}
	if off >= int64(len(r.s)) {
		return 0, io.EOF
	}
	n = copy(b, r.s[off:])
	if n < len(b) {
		err = io.EOF
	}
	return
}

// ReadByte implements the io.ByteReader interface.
// ReadByte方法  实现了 io.ByteReader接口.
// !!! 每次读取一个字节的正常读取方式，会使得当前读取字节的序号加1
func (r *StringReader) ReadByte() (byte, error) {
	r.prevRune = -1 //表示此操作不是ReadRune
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	b := r.s[r.i]
	r.i++
	return b, nil
}

// UnreadByte implements the io.ByteScanner interface.
// UnreadByte实现了io.ByteScanner接口的要求。
// !!! UnreadByte 使字节读取序号退回到上一个字节处。
func (r *StringReader) UnreadByte() error {
	if r.i <= 0 {
		return errors.New("strings.Reader.UnreadByte: at beginning of string")
	}
	r.prevRune = -1 //表示此操作不是ReadRune
	r.i--
	return nil
}

// ReadRune implements the io.RuneReader interface.
// ReadRune 实现了io.RuneReader接口。size表示所返回的字符（rune）所占的字节数。
// UTF-8最多四个字节,所以用int32足以可以表示一个字符，rune就是int32用于表示字符时的别名。
func (r *StringReader) ReadRune() (ch rune, size int, err error) {
	//如果当前读取字节的序号已在字符串合法字节序号边界之外，就返回io.EOF。
	if r.i >= int64(len(r.s)) {
		r.prevRune = -1 //表示此操作不是ReadRune
		return 0, 0, io.EOF
	}
	//将prevRune设置为当前的字节读取序号
	r.prevRune = int(r.i)
	//如果 一个字节(byte)的值小于128（utf8.RuneSelf=128）就是ASC码，在utf-8中只占1个字节
	if c := r.s[r.i]; c < utf8.RuneSelf {
		r.i++
		return rune(c), 1, nil
	}
	//DecodeRuneInString与DecodeRune相似，只不过它的输入是一个字符串。如果s为空，该方法就返回
	//(RuneError, 0),否则，如果编码不正确，就会返回(RuneError, 1)。对于正确的，非空的UTF-8字符
	//串来说，是不可能出现上述两种情况的。RuneError是一个rune字符常量，值为 '\uFFFD'
	//如果字符串不是正确的UTF-8格式、对超过unicod范围的rune进行编码，不是最小可能的UTF-8值，
	//那么编码就会不正确，而没有执行其他校验。
	ch, size = utf8.DecodeRuneInString(r.s[r.i:])
	r.i += int64(size)
	return
}

// UnreadRune implements the io.RuneScanner interface.
//UnreadRunes实现了io.RuneScanner 接口，
//该方法的主要操作就是把当前的读取字节序号回退到前一个字符（rune）所起始的字节序号。
//由于preRune只能记录上一个字符（rune）所起始的字节书号，所以无法再继续回退到再前一个字符（rune）。

func (r *StringReader) UnreadRune() error {
	if r.i <= 0 {
		return errors.New("strings.Reader.UnreadRune: at beginning of string")
	}
	if r.prevRune < 0 {
		return errors.New("strings.Reader.UnreadRune: previous operation was not ReadRune")
	}
	r.i = int64(r.prevRune)
	r.prevRune = -1 //表示此操作不是ReadRune
	return nil
}

// Seek implements the io.Seeker interface.
// Seek 实现了io.Seeker接口。
// !!! Seek是查找的意思。
func (r *StringReader) Seek(offset int64, whence int) (int64, error) {
	r.prevRune = -1
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.i + offset
	case io.SeekEnd:
		abs = int64(len(r.s)) + offset
	default:
		return 0, errors.New("strings.Reader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("strings.Reader.Seek: negative position")
	}
	r.i = abs
	return abs, nil
}

// WriteTo implements the io.WriterTo interface.
func (r *StringReader) WriteTo(w io.Writer) (n int64, err error) {
	r.prevRune = -1
	if r.i >= int64(len(r.s)) {
		return 0, nil
	}
	s := r.s[r.i:]
	// WriteString writes the contents of the string s to w, which accepts a slice of bytes.
	// If w implements StringWriter, its WriteString method is invoked directly.
	// Otherwise, w.Write is called exactly once.
	m, err := io.WriteString(w, s)
	if m > len(s) {
		panic("strings.Reader.WriteTo: invalid WriteString count")
	}
	r.i += int64(m)
	n = int64(m)
	//若实际写入的字符串长度小于待写入字符串的长度，意味着有些字符串没有写入，会返回ErrShortWrite错误。
	if m != len(s) && err == nil {
		err = io.ErrShortWrite
	}
	return
}

// Reset resets the Reader to be reading from s.
func (r *StringReader) Reset(s string) { *r = StringReader{s, 0, -1} }

// NewReader returns a new Reader reading from s.
// It is similar to bytes.NewBufferString but more efficient and read-only.
func NewReader(s string) *StringReader { return &StringReader{s, 0, -1} }
