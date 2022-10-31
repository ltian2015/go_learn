// 这是对GO IO包的翻译。
// Package io provides basic interfaces to I/O primitives.
// Its primary job is to wrap existing implementations of such primitives,
// such as those in package os, into shared public interfaces that
// abstract the functionality, plus some other related primitives.
// io 包提供了I/O原语的基本接口。
// io包的首要工作将诸如在os包中的I/O实现等已经存在I/O原语实现封装入抽象了功能的共享的公开接口，
// 并加上其他的一些相关原语。
// Because these interfaces and primitives wrap lower-level operations with
// various implementations, unless otherwise informed clients should not
// assume they are safe for parallel execution.
// !!! 因为这些接口和原语分那个装了各种不同实现的低级操作，除非是熟知客户端，
// !!! 否则不能假定并行执行这些操作是安全的
package goio

import (
	"errors"
	"sync"
)

// Seek whence values.
const (
	//
	SeekStart = 0 // seek relative to the origin of the file
	//
	SeekCurrent = 1 // seek relative to the current offset
	//
	SeekEnd = 2 // seek relative to the end
)

// ErrShortWrite means that a write accepted fewer bytes than requested
// but failed to return an explicit error.
// ErrShortWrite表示一个写操作实际接收的字节比请求写入的字节少，但有没有返回一个显式的error(panic方式)
var ErrShortWrite = errors.New("short write")

// errInvalidWrite means that a write returned an impossible count.
// errInvalidWrite表示一个写操作返回了一个不可能的（写入）数量。
var errInvalidWrite = errors.New("invalid write result")

// ErrShortBuffer means that a read required a longer buffer than was provided.
// ErrShortBuffer 表示一个读操作需要一个比提供的buffer更长的buffer，也就是现有的buffer长度不够。
var ErrShortBuffer = errors.New("short buffer")

// EOF is the error returned by Read when no more input is available.
// (Read must return EOF itself, not an error wrapping EOF,
// because callers will test for EOF using ==.)
// Functions should return EOF only to signal a graceful end of input.
// If the EOF occurs unexpectedly in a structured data stream,
// the appropriate error is either ErrUnexpectedEOF or some other error
// giving more detail.
// EOF 是一个由读（Read）一类操作在没有更多的输入可以获取时所返回的error。
// !!! (Read操作必须返回EOF自身，而不能是封装了EOF的error，因为调用者会使用==操作来检测EOF)
// 如果在一个结构化的数据流（structured data stream）发生了非预期的EOF，合适的错误值要么是
// ErrUnexcpetedEOF ，要么是其他给出了更多细节的错误值。
var EOF = errors.New("EOF")

// ErrUnexpectedEOF means that EOF was encountered in the
// middle of reading a fixed-size block or data structure.
// ErrUnexpectedEOF表示在读取一个固定长度字节块或者数据结构的中间遇到了EOF错误。
var ErrUnexpectedEOF = errors.New("unexpected EOF")

// ErrNoProgress is returned by some clients of a Reader when
// many calls to Read have failed to return any data or error,
// usually the sign of a broken Reader implementation.
// 当很多个对Read方法的调用没能够返回任何数据或错误，某些Reader的客户端就会返回ErrNoProgress
var ErrNoProgress = errors.New("multiple Read calls return no data or error")

// Reader is the interface that wraps the basic Read method.
// Reader是一个封装了基本Read方法的接口。
//
// Read reads up to len(p) bytes into p. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered. Even if Read
// returns n < len(p), it may use all of p as c space during the call.
// If some data is available but not len(p) bytes, Read conventionally
// returns what is available instead of waiting for more.
// Read读取至多len(p)个字节到切片p中。它返回所读取的字节数量n（0<=n<=len(p)）以及所遇到的任何错误。
// 即便Read方法返回n<len(p),在调用过程中，Read方法也会把p的全部当作暂存空间。如果可以得到一些数据，
// 但是不够len(P)个字节，Read方法照常会返回所得到的数据，而不是等待把切片p填满。
//
// When Read encounters an error or end-of-file condition after
// successfully reading n > 0 bytes, it returns the number of
// bytes read. It may return the (non-nil) error from the same call
// or return the error (and n == 0) from a subsequent call.
// An instance of this general case is that a Reader returning
// a non-zero number of bytes at the end of the input stream may
// return either err == EOF or err == nil. The next Read should
// return 0, EOF.
// 当Read成功读取n>0个字节后，遇到一个错误或者end-of-file情形时，它会返回所读取的字节数。
// 在与返回n>0的同一个调用中，Read方法也可以返回（非空的）error，或者从后续的调用
// 中返回error（并且n==0）.这种抽象场景的一个具体实例就是一个读取器（Reader）在输入流
// 结尾时返回"非0字节数"的同时，也可能会返回err=EOF或err==nil,而在下一个Read操作时就会返回0，EOF。
// Callers should always process the n > 0 bytes returned before
// considering the error err. Doing so correctly handles I/O errors
// that happen after reading some bytes and also both of the
// allowed EOF behaviors.
// 调用者在考虑error errr之前，应当总会处理所返回的n>0个字节。用这样正确地做法处理发发生
// 在读取一些字节之后的I/O错误，以及所允许的全部两种EOF行为所带来的I/O错误。
// Implementations of Read are discouraged from returning a
// zero byte count with a nil error, except when len(p) == 0.
// Callers should treat a return of 0 and nil as indicating that
// nothing happened; in particular it does not indicate EOF.
//!!! 除非len(P)==0,否则，不鼓励Read的实现返回0个字节和nil error。调用者应当将返回0个字
//!!! 节与nil error 这种情况当作什么也没发生，尤其是没有指明EOF时，更应如此。

// Implementations must not retain p.
// !!! Reader的实现者一定不要持有p。
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Writer is the interface that wraps the basic Write method.
// Writer是一个封装了Write方法的接口。
// Write writes len(p) bytes from p to the underlying data stream.
// It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
// Write must return a non-nil error if it returns n < len(p).
// Write must not modify the slice data, even temporarily.
// Write 方法从p中拿出len(p)个字节写入到底层的数据流之中。它会返回自p中写入到底层
// 数据流中的字节数量（0<=n<=len(p)）,和遇到的，令写入操作提前终止的错误。
// 如果返回的n<len(p),那么Write方法必须返回一个非空（non-nil）的错误(error).
// !!! Write方法一定不能更改切片p中的数据，即使是临时更改也不行。
// Implementations must not retain p.
// !!! Writer的实现一定不要持有p
type Writer interface {
	Write(p []byte) (n int, err error)
}

// Closer is the interface that wraps the basic Close method.
// Closer是一个封装了基本Close方法的接口。
// The behavior of Close after the first call is undefined.
// Specific implementations may document their own behavior.
// 第一次调用Close之后在调用Close方法，其行为是未定义，也就是未知的。
// 特定的实现会说命名它们自己的Close行为。
type Closer interface {
	Close() error
}

// Seeker is the interface that wraps the basic Seek method.
// Seeker 是对基本的Seek方法的封装。
// Seek sets the offset for the next Read or Write to offset,
// interpreted according to whence:
// SeekStart means relative to the start of the file,
// SeekCurrent means relative to the current offset, and
// SeekEnd means relative to the end
// (for example, offset = -2 specifies t he penultimate byte of the file).
// Seek returns the new offset relative to the start of the
// file or an error, if any.
// !!! Seek（搜寻）方法为下一次读（Read）或写（Write）操作设置了偏移（offset），
// 偏移的含义依据“根源（whence）”的含义进行解释：
// 当Whence 值为SeekStart (常量SeekStart==0)s时，offset表示相对于文件开始。
// 当Whence 值为SeekCurrent(常量SeekCurrent==1)，offset表示相对于当前位置的偏移量，而
// 当Whence 值为SeekEnd（常量SeekEnd==2）时，offset表示表示相对于文件结尾的偏移量。(例如，
// 当offset=-2 时，表示的是文件倒数第二个字节的位置)
// !!! Seek方法返回相对于文件开始的新的偏移量，或者错误，如果有错误的话。
// Seeking to an offset before the start of the file is an error.
// Seeking to any positive offset may be allowed, but if the new offset exceeds
// the size of the underlying object the behavior of subsequent I/O operations
// is implementation-dependent.
// Seeking（搜寻）指定文件开始之前的偏移（offset）是一个错误（error）。
// Seeking（搜寻）到任何一个正偏移都是允许的，但是如果这个新的偏移超过了底层对象个数的大小，
// 那么后续的I/O操作的行为就取决各自的实现了。
// !!! Seek方法本身并不接口使得Reader或Writer可以以非连续的方式进行流的读或写操作。通常
// !!! 用于通信协议的写入与解析，这类文件的字节之间有相对位置的要求。
type Seeker interface {
	Seek(offset int64, whence int) (int64, error)
}

// ReadWriter is the interface that groups the basic Read and Write methods.
// ReadWriter是组合了基本的Read和Write方法的接口。
type ReadWriter interface {
	Reader
	Writer
}

// ReadCloser is the interface that groups the basic Read and Close methods.
// ReadCloser 是组合了基本的Read和Close方法的接口。
type ReadCloser interface {
	Reader
	Closer
}

// WriteCloser is the interface that groups the basic Write and Close methods.
// WirteCloser是组合了基本的Write和Close方法的接口
type WriteCloser interface {
	Writer
	Closer
}

// ReadWriteCloser is the interface that groups the basic Read, Write and Close methods.
// ReadWriteCloser 是组合了基本的Read，Write，Close方法的接口。
type ReadWriteCloser interface {
	Reader
	Writer
	Closer
}

// ReadSeek 是组合了基本的Read和Seek方法的接口
// ReadSeeker is the interface that groups the basic Read and Seek methods.
type ReadSeeker interface {
	Reader
	Seeker
}

// ReadSeekCloser is the interface that groups the basic Read, Seek and Close
// methods.
// ReadSeekCloser是组合了基本的Read，Seek和Close方法的接口。
type ReadSeekCloser interface {
	Reader
	Seeker
	Closer
}

// WriteSeeker is the interface that groups the basic Write and Seek methods.
// WriteSeeker是组合了基本的Write和Seek方法的接口。
type WriteSeeker interface {
	Writer
	Seeker
}

// ReadWriteSeeker is the interface that groups the basic Read, Write and Seek methods.
// ReadWriteSeeker 是组合了基本的Read，Write和Seek方法的接口。
type ReadWriteSeeker interface {
	Reader
	Writer
	Seeker
}

// ReaderFrom is the interface that wraps the ReadFrom method.
//
// ReadFrom reads data from r until EOF or error.
// The return value n is the number of bytes read.
// Any error except EOF encountered during the read is also returned.
//
// The Copy function uses ReaderFrom if available.

// ReadFrom 是一个封装了ReadFrom方法的接口。
// ReadFrom 从读取器r读取数据到内存中，直到遇到EOF或者error。
// 返回值是所读取的字节数。
// 除了EOF之外，在读取期间遇到的任何error都会被返回。
// 如果可能，Copy函数使用ReaderFrom接口
// !!! ReaderFrom 接口的实现内部在内存中应持有从Reader中所读取的数据。
type ReaderFrom interface {
	ReadFrom(r Reader) (n int64, err error)
}

// WriterTo is the interface that wraps the WriteTo method.
//
// WriteTo writes data to w until there's no more data to write or
// when an error occurs. The return value n is the number of bytes
// written. Any error encountered during the write is also returned.
//
// The Copy function uses WriterTo if available.

// WriteTo是一个封装了WriteTo方法的接口。
// WriteTo 将数据内存数据写入到写入器中，直到没有数据可写，或者遇到了错误发生。
// 返回值n是写入的字节数量。在写入的过程中，遇到任何错误，都会被返回。
// 如果可能，Copy函数使用WriteTo接口
// !!! WriterTo接口的实现内部应持有写入到w中的内存数据。
type WriterTo interface {
	WriteTo(w Writer) (n int64, err error)
}

// ReaderAt is the interface that wraps the basic ReadAt method.
// ReaderAt 是一个封装了基本的ReadAt方法的接口。
// ReadAt reads len(p) bytes into p starting at offset off in the
// underlying input source. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered.
// ReadAt 从偏移处（off参数）开始，从底层的输入源中读取len（p）个字节到p中。
// 它会返回所所读取的字节数（0 <= n <= len(p)）,以及所遇到任何错误。
// When ReadAt returns n < len(p), it returns a non-nil error
// explaining why more bytes were not returned. In this respect,
// ReadAt is stricter than Read.
// !!! 当ReadAt返回n<len(p) 时，它一定会返回一个非空（non-nil）error来解释为什么
// !!! 没有返回更多的字节。在这一点，ReadAt要比Read更加严格。
// !!! （在这种情况下Read接口不一定会返回EOF ,err!=nil=,而是在下一个read调用时返回(0,err)
// Even if ReadAt returns n < len(p), it may use all of p as scratch
// space during the call. If some data is available but not len(p) bytes,
// ReadAt blocks until either all the data is available or an error occurs.
// In this respect ReadAt is different from Read.
// 尽管ReadAt返回n<len(p),在调用过程中，它可以使用P的所有空间作为暂存空间。
// !!! 如果可以得到一些数据，但是不是len(p)个字节。ReadAt会阻塞，直到得到所有
// !!! 的数据，或者遇到一个错误。在这一点上，ReadAt与Read完全不同。
// If the n = len(p) bytes returned by ReadAt are at the end of the
// input source, ReadAt may return either err == EOF or err == nil.
// 如果在输入流结束的时候，ReadAt返回 n=len（p）（也就是输入流的长度与p的长度正好相等），
// ReadAt可以返回err=EOF或者err=nil。
// If ReadAt is reading from an input source with a seek offset,
// ReadAt should not affect nor be affected by the underlying
// seek offset.
// !!! 如果ReadAt正在从读取一个带有seek offset的输入源，
// !!! ReadAt不应影响底层的seek offset，也不应被底层的seek offset 所影响。
// Clients of ReadAt can execute parallel ReadAt calls on the
// same input source.
// !!! ReadAt客户端可以在同一个输入源上并行执行多个ReadAt调用。
// Implementations must not retain p.
// !!! ReadAt接口的实现一定不要持有p。
type ReaderAt interface {
	ReadAt(p []byte, off int64) (n int, err error)
}

// WriterAt is the interface that wraps the basic WriteAt method.
// WriterAt 是封装了基本的WriteAt方法的接口。

// WriteAt writes len(p) bytes from p to the underlying data stream
// at offset off. It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
// WriteAt must return a non-nil error if it returns.

// WriteAt在偏移off处，把p中的len(p)个字节到底层的数据流中。它返回从p所写入的
// 字节数(0 <= n <= len(p)),以及所遇到的任何引起write操作提前终止的错误。
// 如果0 <= n <= len(p)，那么WriteAt必须返回非空的error。

// If WriteAt is writing to a destination with a seek offset,
// WriteAt should not affect nor be affected by the underlying
// seek offset.

// 如果WriteAt正在向一个带有seek offset的目标写入字节，
// WriteAt本应影响底层数据流的seek offset，也不应被底层数据流的seek offset所影响。

// Clients of WriteAt can execute parallel WriteAt calls on the same
// destination if the ranges do not overlap.

// !!! 如果写入的范围不发生重叠，WriteAt的客户端可以在同一个写入目标上并行地执行多个WriteAt调用。
// Implementations must not retain p.
// !!!  WriteAt的实现一定不能持有p。
type WriterAt interface {
	WriteAt(p []byte, off int64) (n int, err error)
}

// ByteReader is the interface that wraps the ReadByte method.
//
// ReadByte reads and returns the next byte from the input or
// any error encountered. If ReadByte returns an error, no input
// byte was consumed, and the returned byte value is undefined.
//
// ReadByte provides an efficient interface for byte-at-time
// processing. A Reader that does not implement  ByteReader
// can be wrapped using bufio.NewReader to add this method.
type ByteReader interface {
	ReadByte() (byte, error)
}

// ByteScanner is the interface that adds the UnreadByte method to the
// basic ReadByte method.
//
// UnreadByte causes the next call to ReadByte to return the last byte read.
// If the last operation was not a successful call to ReadByte, UnreadByte may
// return an error, unread the last byte read (or the byte prior to the
// last-unread byte), or (in implementations that support the Seeker interface)
// seek to one byte before the current offset.
type ByteScanner interface {
	ByteReader
	UnreadByte() error
}

// ByteWriter is the interface that wraps the WriteByte method.
type ByteWriter interface {
	WriteByte(c byte) error
}

// RuneReader is the interface that wraps the ReadRune method.
//
// ReadRune reads a single encoded Unicode character
// and returns the rune and its size in bytes. If no character is
// available, err will be set.
type RuneReader interface {
	ReadRune() (r rune, size int, err error)
}

// RuneScanner is the interface that adds the UnreadRune method to the
// basic ReadRune method.
//
// UnreadRune causes the next call to ReadRune to return the last rune read.
// If the last operation was not a successful call to ReadRune, UnreadRune may
// return an error, unread the last rune read (or the rune prior to the
// last-unread rune), or (in implementations that support the Seeker interface)
// seek to the start of the rune before the current offset.
type RuneScanner interface {
	RuneReader
	UnreadRune() error
}

// StringWriter is the interface that wraps the WriteString method.
type StringWriter interface {
	WriteString(s string) (n int, err error)
}

// WriteString writes the contents of the string s to w, which accepts a slice of bytes.
// If w implements StringWriter, its WriteString method is invoked directly.
// Otherwise, w.Write is called exactly once.
func WriteString(w Writer, s string) (n int, err error) {
	if sw, ok := w.(StringWriter); ok {
		return sw.WriteString(s)
	}
	return w.Write([]byte(s))
}

// ReadAtLeast reads from r into buf until it has read at least min bytes.
// It returns the number of bytes copied and an error if fewer bytes were read.
// The error is EOF only if no bytes were read.
// If an EOF happens after reading fewer than min bytes,
// ReadAtLeast returns ErrUnexpectedEOF.
// If min is greater than the length of buf, ReadAtLeast returns ErrShortBuffer.
// On return, n >= min if and only if err == nil.
// If r returns an error having read at least min bytes, the error is dropped.
func ReadAtLeast(r Reader, buf []byte, min int) (n int, err error) {
	if len(buf) < min {
		return 0, ErrShortBuffer
	}
	for n < min && err == nil {
		var nn int
		nn, err = r.Read(buf[n:])
		n += nn
	}
	if n >= min {
		err = nil
	} else if n > 0 && err == EOF {
		err = ErrUnexpectedEOF
	}
	return
}

// ReadFull reads exactly len(buf) bytes from r into buf.
// It returns the number of bytes copied and an error if fewer bytes were read.
// The error is EOF only if no bytes were read.
// If an EOF happens after reading some but not all the bytes,
// ReadFull returns ErrUnexpectedEOF.
// On return, n == len(buf) if and only if err == nil.
// If r returns an error having read at least len(buf) bytes, the error is dropped.
func ReadFull(r Reader, buf []byte) (n int, err error) {
	return ReadAtLeast(r, buf, len(buf))
}

// CopyN copies n bytes (or until an error) from src to dst.
// It returns the number of bytes copied and the earliest
// error encountered while copying.
// On return, written == n if and only if err == nil.
//
// If dst implements the ReaderFrom interface,
// the copy is implemented using it.
func CopyN(dst Writer, src Reader, n int64) (written int64, err error) {
	written, err = Copy(dst, LimitReader(src, n))
	if written == n {
		return n, nil
	}
	if written < n && err == nil {
		// src stopped early; must have been EOF.
		err = EOF
	}
	return
}

// Copy copies from src to dst until either EOF is reached
// on src or an error occurs. It returns the number of bytes
// copied and the first error encountered while copying, if any.
//
// A successful Copy returns err == nil, not err == EOF.
// Because Copy is defined to read from src until EOF, it does
// not treat an EOF from Read as an error to be reported.
//
// If src implements the WriterTo interface,
// the copy is implemented by calling src.WriteTo(dst).
// Otherwise, if dst implements the ReaderFrom interface,
// the copy is implemented by calling dst.ReadFrom(src).
func Copy(dst Writer, src Reader) (written int64, err error) {
	return copyBuffer(dst, src, nil)
}

// CopyBuffer is identical to Copy except that it stages through the
// provided buffer (if one is required) rather than allocating a
// temporary one. If buf is nil, one is allocated; otherwise if it has
// zero length, CopyBuffer panics.
//
// If either src implements WriterTo or dst implements ReaderFrom,
// buf will not be used to perform the copy.
func CopyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) {
	if buf != nil && len(buf) == 0 {
		panic("empty buffer in CopyBuffer")
	}
	return copyBuffer(dst, src, buf)
}

// copyBuffer is the actual implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.
func copyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	if wt, ok := src.(WriterTo); ok {
		return wt.WriteTo(dst)
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := dst.(ReaderFrom); ok {
		return rt.ReadFrom(src)
	}
	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

// LimitReader returns a Reader that reads from r
// but stops with EOF after n bytes.
// The underlying implementation is a *LimitedReader.
func LimitReader(r Reader, n int64) Reader { return &LimitedReader{r, n} }

// A LimitedReader reads from R but limits the amount of
// data returned to just N bytes. Each call to Read
// updates N to reflect the new amount remaining.
// Read returns EOF when N <= 0 or when the underlying R returns EOF.
type LimitedReader struct {
	R Reader // underlying reader
	N int64  // max bytes remaining
}

func (l *LimitedReader) Read(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, EOF
	}
	if int64(len(p)) > l.N {
		p = p[0:l.N]
	}
	n, err = l.R.Read(p)
	l.N -= int64(n)
	return
}

// NewSectionReader returns a SectionReader that reads from r
// starting at offset off and stops with EOF after n bytes.
func NewSectionReader(r ReaderAt, off int64, n int64) *SectionReader {
	var remaining int64
	const maxint64 = 1<<63 - 1
	if off <= maxint64-n {
		remaining = n + off
	} else {
		// Overflow, with no way to return error.
		// Assume we can read up to an offset of 1<<63 - 1.
		remaining = maxint64
	}
	return &SectionReader{r, off, off, remaining}
}

// SectionReader implements Read, Seek, and ReadAt on a section
// of an underlying ReaderAt.
type SectionReader struct {
	r     ReaderAt
	base  int64
	off   int64
	limit int64
}

func (s *SectionReader) Read(p []byte) (n int, err error) {
	if s.off >= s.limit {
		return 0, EOF
	}
	if max := s.limit - s.off; int64(len(p)) > max {
		p = p[0:max]
	}
	n, err = s.r.ReadAt(p, s.off)
	s.off += int64(n)
	return
}

var errWhence = errors.New("Seek: invalid whence")
var errOffset = errors.New("Seek: invalid offset")

func (s *SectionReader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	default:
		return 0, errWhence
	case SeekStart:
		offset += s.base
	case SeekCurrent:
		offset += s.off
	case SeekEnd:
		offset += s.limit
	}
	if offset < s.base {
		return 0, errOffset
	}
	s.off = offset
	return offset - s.base, nil
}

func (s *SectionReader) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 || off >= s.limit-s.base {
		return 0, EOF
	}
	off += s.base
	if max := s.limit - off; int64(len(p)) > max {
		p = p[0:max]
		n, err = s.r.ReadAt(p, off)
		if err == nil {
			err = EOF
		}
		return n, err
	}
	return s.r.ReadAt(p, off)
}

// Size returns the size of the section in bytes.
func (s *SectionReader) Size() int64 { return s.limit - s.base }

// TeeReader returns a Reader that writes to w what it reads from r.
// All reads from r performed through it are matched with
// corresponding writes to w. There is no internal buffering -
// the write must complete before the read completes.
// Any error encountered while writing is reported as a read error.
func TeeReader(r Reader, w Writer) Reader {
	return &teeReader{r, w}
}

type teeReader struct {
	r Reader
	w Writer
}

func (t *teeReader) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 {
		if n, err := t.w.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return
}

// Discard is a Writer on which all Write calls succeed
// without doing anything.
var Discard Writer = discard{}

type discard struct{}

// discard implements ReaderFrom as an optimization so Copy to
// io.Discard can avoid doing unnecessary work.
var _ ReaderFrom = discard{}

func (discard) Write(p []byte) (int, error) {
	return len(p), nil
}

func (discard) WriteString(s string) (int, error) {
	return len(s), nil
}

var blackHolePool = sync.Pool{
	New: func() any {
		b := make([]byte, 8192)
		return &b
	},
}

func (discard) ReadFrom(r Reader) (n int64, err error) {
	bufp := blackHolePool.Get().(*[]byte)
	readSize := 0
	for {
		readSize, err = r.Read(*bufp)
		n += int64(readSize)
		if err != nil {
			blackHolePool.Put(bufp)
			if err == EOF {
				return n, nil
			}
			return
		}
	}
}

// NopCloser returns a ReadCloser with a no-op Close method wrapping
// the provided Reader r.
// If r implements WriterTo, the returned ReadCloser will implement WriterTo
// by forwarding calls to r.
func NopCloser(r Reader) ReadCloser {
	if _, ok := r.(WriterTo); ok {
		return nopCloserWriterTo{r}
	}
	return nopCloser{r}
}

type nopCloser struct {
	Reader
}

func (nopCloser) Close() error { return nil }

type nopCloserWriterTo struct {
	Reader
}

func (nopCloserWriterTo) Close() error { return nil }

func (c nopCloserWriterTo) WriteTo(w Writer) (n int64, err error) {
	return c.Reader.(WriterTo).WriteTo(w)
}

// ReadAll reads from r until an error or EOF and returns the data it read.
// A successful call returns err == nil, not err == EOF. Because ReadAll is
// defined to read from src until EOF, it does not treat an EOF from Read
// as an error to be reported.
func ReadAll(r Reader) ([]byte, error) {
	b := make([]byte, 0, 512)
	for {
		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == EOF {
				err = nil
			}
			return b, err
		}
	}
}
