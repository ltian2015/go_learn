// !!! 这是对GO IO包的翻译。
// !!! GO IO 主要负责把程序员当前所开发程序中的内存中的数据通过Writer写入到底层流中，
// !!! 或者将底层流中的数据通过Reader读入到内存中进行处理。写入与读取的数据主要存在
// !!! 两种形式:人不可理解的二进制字节(byte),以及人可以理解的字符（rune）及其衍生的字符串（string）
// !!! 所谓读写都是针对流与程序员所控制和处理的内存的关系而定，把流中数据拿到内存中叫做"读(Read)"，
// !!! 把内存中数据拿到流中叫做"写(Write)"。所谓“流”，就是一切可以吞吐数据的东西，任何具有读写数据接口的
// !!! 的组件都是“流”。而GO I/O包就是定义了流的标准接口、工具方法和特定能力的具体实现。

// !!! GO IO包 主要的定义了“流组件”的标准接口，基于这些接口的工具方法，以及这些接口的一些特殊实现，主要是：

// !!!  ------------流组件的标准接口如下：
// !!! Reader: 基本的“批量字节读取器”
// !!! Writer: 基本的“批量字节写入器”
// !!! Closer: 读写器关闭的控制器
// !!! Seeker: 使读、写器前往到特定位置的“读写位置控制器”，代表可反复读取、可跳跃式读取的读取器。
// !!! ReaderFrom: 读取器客户端，通过给定的读取器来读取数据。
// !!! WriteTo: 写入器的客户端，通过给定的写入器来写入数据，
// !!! ReaderAt: 在指定偏移处读取批量数据到内存，可用于大量数据并行、分段读取。
// !!! WriteAt: 在指定偏移处写入批量数据到底层流，可用于大量数据并行、分段写入。
// !!! ByteReader:一次只读取一个字节的“逐个字节读取器”。
// !!! ByteWriter: 一次只向底层流写入一个字节的“逐个字节写入器”。
// !!! ByteScanner: 带有回退到上一个读取位置功能的“逐个字节读取器”。
// !!! RuneReader:  一次读取一个utf-8字符的“逐个字符读取器”。
// !!! RuneScanner: 带有退回到当一个utf-8字符所在位置的“逐个字符读取器”
// !!! StringWriter:  将utf8格式字符串写入其他流中的“字符串写入器”
// !!!----------操作流组件的便利的工具方法如下：-------------------------
// !!! WriteString(): 该方法把给定的字符串s写入到给定的“字节写入器w”中。
// !!! ReadAtLeast(): 该方法使用给定的“字节读取器r”读取给定的至少min个字节的数据到给定的字节切片内存中。
// !!! ReadFull(): 该方法通过给定的“字节读取器r”把给定的内存缓存buf读满。
// !!! CopyN(): 该方法将给定的“源（src）-字节读取器”，拷贝给定的n个字节到“目标（dst）——字节写入器”。
// !!! Copy(): 该方法将给定的“源（src）-字节读取器”中所有的数据拷贝“目标（dst）——字节写入器”中。除非全部读取（遇到EOF）或者遇到读写错误。
// !!! CopyBuffer():该通过给定缓存将读取器中读取的数据写入到写入器中。通过设定缓存的大小使用程序员可以优化读、写速度。
// !!! ReadAll(): 该方法试图“返回”给定的字节读取器中所有的有效数据。
// !!! ------------特定的具体流组件实现，及其构建方法和全局变量-------------
// !!! LimitedReader: 从给定的“字节读取器（Reader）”中读取返回的“总计数据量”限制为正好为N个字节，(而无论读几次)。
// !!! LimitReader(): 该方法用给定的“字节读取器”和限制读取的“总计数据量”，创建一个限制总计读取数据数量的字节读取器。
// !!! SectionReader: 是在ReaderAt上构建更加高级的“分段读取器”。一旦分段确定，可以支持多次普通读、再分段读、反复、跳跃等多读取方式。
// !!! NewSectionReader(): 用给定ReaderAt读取器、起始位置、偏移长度，分段的字节长度来构建一个“分段读取器”。
// !!! teeReader:  “T型的三通读取器”是io包私有的读取器类型，其Read方法实现了一次读取，同时把数据送入内存缓存和写入器（另一个流）。
// !!! TeeReader():该方法通过给定的“写入器”参数创建“T型的三通读取器”。
// !!!  discard: 用于抛弃所有写入数据的字节“特殊写入器”，并实现了ReaderFrom接口，优化该写入器作为Copy操作时的性能。
// !!!  var Discard: 是discard类型的全局唯一的实例。discard
// !!! nopCloser:
// !!! NopCloser():
// !!! nopCloserWriterTo:
//
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

// !!! 因为这些接口和原语封装了各种不同IO实现方式的低级操作，除非是熟知客户端，
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
// !!! 底层流有多少数据是未知的，因此，Reader只有尽量按照给定的内存缓存的长度进行数据的读取。
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
// !!! 底层流中能够写入多少字节是未知的，所以Writer只能尽量把内存缓存中的所有数据写入到底层流中。
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
// !!! Seek（前往）方法为下一次读（Read）或写（Write）操作设置了偏移（offset），
// 偏移的含义依据“根源（whence）”的含义进行解释：
// 当Whence 值为SeekStart (常量SeekStart==0)s时，offset表示相对于文件开始。
// 当Whence 值为SeekCurrent(常量SeekCurrent==1)，offset表示相对于当前位置的偏移量，而
// 当Whence 值为SeekEnd（常量SeekEnd==2）时，offset表示表示相对于文件结尾的偏移量。(例如，
// 当offset=-2 时，表示的是文件倒数第二个字节的位置)
// !!! 注意，当Whence为SeekEnd时，有效的offset应为负数，否则读取的字节数会为0.
// !!! Seek方法返回相对于文件开始的新的偏移量，或者错误，如果有错误的话。
// Seeking to an offset before the start of the file is an error.
// Seeking to any positive offset may be allowed, but if the new offset exceeds
// the size of the underlying object the behavior of subsequent I/O operations
// is implementation-dependent.
// Seeking（前往）指定文件开始之前的偏移（offset）是一个错误（error）。
// Seeking（前往）到任何一个正偏移都是允许的，但是如果这个新的偏移超过了底层对象个数的大小，
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
// !!! 因为使用外部传递过来的Reader写入数据，ReadFrom 定义的是Reader客户端程序的接口。
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
// !!! 因为使用外部传递过来的Reader写入数据，WriterTo 接口定义的是Writer的客户端程序。
type WriterTo interface {
	WriteTo(w Writer) (n int64, err error)
}

// ReaderAt is the interface that wraps the basic ReadAt method.

// ReaderAt 是一个封装了基本的ReadAt方法的接口。

// ReadAt reads len(p) bytes into p starting at offset off in the
// underlying input source. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered.

// !!! ReadAt 从偏移处（off参数）开始，从底层的输入源中读取len（p）个字节到p中。
// !!! 它会返回所所读取的字节数（0 <= n <= len(p)）,以及所遇到任何错误。

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
// !!! ReadAt从某偏移处读取指定长度的字节，可以用于并行、分段读取大量定长数据。
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
// !!!  WriteAt可用于并行、分段写入大量数据。
type WriterAt interface {
	WriteAt(p []byte, off int64) (n int, err error)
}

// ByteReader is the interface that wraps the ReadByte method.

// ByteREader 是一个封装le ReadByte方法的接口。

// ReadByte reads and returns the next byte from the input or
// any error encountered. If ReadByte returns an error, no input
// byte was consumed, and the returned byte value is undefined.

// ReadByte 从输入流汇中读取并返回下一个字节或者所遇到的任何错误。如果
// ReadByte 返回一个错误，那就没有输入字节被消费，并且返回的字节值将是未定义的（无效的）.

// ReadByte provides an efficient interface for byte-at-time
// processing. A Reader that does not implement  ByteReader
// can be wrapped using bufio.NewReader to add this method.

// !!! ReadByte 为“一次一个字节的处理”提供了一种有效的接口。没有实现
// ByteReader的Reader可以通过使用bufio.NewReader来增加这个方法，
// 从而对原Reader进行新的包装。
type ByteReader interface {
	ReadByte() (byte, error)
}

// ByteScanner is the interface that adds the UnreadByte method to the
// basic ReadByte method.
//
// ByteScanner 是在基本能的ReadByte方法上增加了一个UnreadByte方法的接口。

// UnreadByte causes the next call to ReadByte to return the last byte read.
// If the last operation was not a successful call to ReadByte, UnreadByte may
// return an error, unread the last byte read (or the byte prior to the
// last-unread byte), or (in implementations that support the Seeker interface)
// seek to one byte before the current offset.
// UnreadByte 导致下一个ReadByte调用返回上一次所读的字节（byte）。

// 如果上一个调用ReadByte的操作没有成功，UnreadByte可以返回一个错误（error），unread上一
// 个字节的read(或者是上一个未读字节的前一个字节)，或者（在支持Seeker接口的
// 实现中）前往当前偏移的上一个字节 。
type ByteScanner interface {
	ByteReader
	UnreadByte() error
}

// ByteWriter is the interface that wraps the WriteByte method.
// ByteWriter 是一个封装了WriteByte方法的接口。
// !!! 与ByteReader一样，主要用于“一次一个字节处理”的工作场景。
type ByteWriter interface {
	WriteByte(c byte) error
}

// RuneReader is the interface that wraps the ReadRune method.

// RuneReader是封装了ReadRune方法的接口。

// ReadRune reads a single encoded Unicode character
// and returns the rune and its size in bytes. If no character is
// available, err will be set.
//
// ReadRune读取了一个单一的（utf-8）编码后的Unicode字符，并且返回该字符所占用的
// 字节数。如果没有读取到字符，那就必须设置err。
type RuneReader interface {
	ReadRune() (r rune, size int, err error)
}

// RuneScanner is the interface that adds the UnreadRune method to the
// basic ReadRune method.

// RuneScanner是在基本的ReadRune方法上又增加了UnreadRune方法的接口。

// UnreadRune causes the next call to ReadRune to return the last rune read.
// If the last operation was not a successful call to ReadRune, UnreadRune may
// return an error, unread the last rune read (or the rune prior to the
// last-unread rune), or (in implementations that support the Seeker interface)
// seek to the start of the rune before the current offset.

// UnreadRune 导致下一次对ReadRune的调用返回上一次所读取的rune。如果上一次ReadRune调用没有
// 成功，UnreadRune可以返回一个错误，unread上一个rune read（或者上一个未读的rune的前一个rune),
// 或者（在实现了支持Seeker接口的实现中）前往（seek）当前偏移的前一个rune。
type RuneScanner interface {
	RuneReader
	UnreadRune() error
}

// StringWriter is the interface that wraps the WriteString method.
// StringWriter是一个封装了WriteString方法的接口。
// !!! StringWriter把内存中的字符串写入底层流中。
// !!! 返回的n应该是“字节数”
// !!! StringWriter与Writer类似，只不过写入的不是字节切片，而是字符串。
type StringWriter interface {
	WriteString(s string) (n int, err error)
}

// WriteString writes the contents of the string s to w, which accepts a slice of bytes.
// If w implements StringWriter, its WriteString method is invoked directly.
// Otherwise, w.Write is called exactly once.
// WriteString方法把字符串s写入到写入器w中，写入器w可以接受字节切片。
// 如果w实现了StringWriter，则会直接调用WriteString方法。
// 否则，w.Write会被调用一次。
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
//
// ReadAtLeast 从r中读取数据到内存buf中，直到它已经至少读取了min个字节。它会返回被拷贝的字节数,
// 并且，如果读取的数据量不足(达不到min).它会返回一个错误。
// 只有当没有字节可读的时候，才会返回EOF 这个错误。如果所读读的数据少于min个字节时发生了EOF，
// ReadAtLeast就会返回ErrUnexpectedEOF。
//
// If min is greater than the length of buf, ReadAtLeast returns ErrShortBuffer.
// On return, n >= min if and only if err == nil.

// 如果min比buf的长度大（buf容量不足），ReadAtLeast就会返回ErrShortBuffer。
// 当返回的时候,并且仅当err==nil时，n>=min。

// If r returns an error having read at least min bytes, the error is dropped.

// 如果在读取了至少min个字节后，r返回了一个错误，该错误会被抛弃。
// !!! 这是IO包中提供的一个工具方法。
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

// ReadFull通过读取器r读取恰好len(buf)个字节到内存缓存buf中。
// 它返回所拷贝的字节数，以及如果字节读取不足时的错误。
// 只有当没有字节被读取时，error才会是EOF。
// If an EOF happens after reading some but not all the bytes,
// ReadFull returns ErrUnexpectedEOF.

// 如果当EOF发生在读取一部分字节但不是全部字节的时候，ReadFull返回ErrUnexpectedEOF

// On return, n == len(buf) if and only if err == nil.
// If r returns an error having read at least len(buf) bytes, the error is dropped.
// 当返回时，只有当err==nil时，n==len(buf)。
// 如果在r已经读取了至少len（buf）个字节后时返回了一个错误，这个错误会被丢弃。
// !!! 这是IO包中提供的一个工具方法。
func ReadFull(r Reader, buf []byte) (n int, err error) {
	return ReadAtLeast(r, buf, len(buf))
}

// CopyN copies n bytes (or until an error) from src to dst.
// It returns the number of bytes copied and the earliest
// error encountered while copying.
// On return, written == n if and only if err == nil.
//
// CopyN 从源（src）拷贝n个字节到目标（dst）。
// 它返回所拷贝的字节数量，以及在拷贝过程中所遇到的最早的错误。
// 在返回后，只有err==nil，才会有wrtten==n 。

// If dst implements the ReaderFrom interface,
// the copy is implemented using it.

// 如果dst实现了ReaderFrom接口，那么Copy就用它来实现。
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

// Copy方法从src 拷贝字节到dst，直到src遇到了EOF或者发生了错误。它返回
// 所拷贝的字节数量，以及在拷贝过程中可能遇到的错误。

// A successful Copy returns err == nil, not err == EOF.
// Because Copy is defined to read from src until EOF, it does
// not treat an EOF from Read as an error to be reported.
//
// !!! 一个成功的Copy调用，返回err=nil，而不是err==EOF。
// 因为Copy被定义为直到src遇到EOF，所以它不会将Read遇到的EOF当做错误进行报告。

// If src implements the WriterTo interface,
// the copy is implemented by calling src.WriteTo(dst).
// Otherwise, if dst implements the ReaderFrom interface,
// the copy is implemented by calling dst.ReadFrom(src).

// 如果src实现了WriteTo接口，拷贝就会用src.WriteTo(dst)来实现。
// 否则，如果dst实现了ReaderFrom接口，拷贝就会用dst.ReadFrom(src)来实现。

func Copy(dst Writer, src Reader) (written int64, err error) {
	return copyBuffer(dst, src, nil)
}

// CopyBuffer is identical to Copy except that it stages through the
// provided buffer (if one is required) rather than allocating a
// temporary one. If buf is nil, one is allocated; otherwise if it has
// zero length, CopyBuffer panics.
//
// CopyBuffer除了使用一个提供的buffer而不是使用一个临时开辟的buffer之外，与Copy
// 方法没有区别。如果提供的buf是nil，就会开辟一个；否则，如果buf的长度为0，CopyBuffer就会panic。

// If either src implements WriterTo or dst implements ReaderFrom,
// buf will not be used to perform the copy.
//
// 如果src实现了WriterTo或者dst实现了ReaderFrom，那么在拷贝过程中就不会使用buf。
// !!! CopyBuffer允许给定的buf为nil，但不允许给定的buf不为nil但是长度为0.
// !!! CopyBuffer通过给定缓存将读取器中读取的数据写入到写入器中。
// !!! CopyBuffer允许开发者通过设定缓存的大小来优化读、写的匹配速度与内存开销。
func CopyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) {
	if buf != nil && len(buf) == 0 {
		panic("empty buffer in CopyBuffer")
	}
	return copyBuffer(dst, src, buf)
}

// copyBuffer is the actual implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.

// copyBuffer是Copy和CopyBuffer实际实现。
//
//	!!! 如果buf是nil，就会开辟一个。
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

// LimitReader返回一个Reader，该Reader从r中读取字节，但是会在读取n个字节后，
// 以EOF结束。这个Reader的底层的实现是*LimitedReader。
// !!! LimitReader方法创建了一个在给定的“字节读取器”基础增加了
// !!! 限制“总计读取字节数量”到（一个或多个）内存缓存buf中的“字节读取器”
func LimitReader(r Reader, n int64) Reader { return &LimitedReader{r, n} }

// A LimitedReader reads from R but limits the amount of
// data returned to just N bytes. Each call to Read
// updates N to reflect the new amount remaining.
// Read returns EOF when N <= 0 or when the underlying R returns EOF.

// LimitedReader从给定的“字节读取器”中读取返回的“总计数据量”限制为正好为N个字节，(而无论读几次)。
// 每次对Read方法的调用都会更新N，以反应出剩余未读字节数量的最新值。
// 当N<=0时或者底层R返回EOF时，LimitedReader的Read方法会返回EOF。
type LimitedReader struct {
	R Reader // underlying reader 底层的读取器
	N int64  // max bytes remaining 最大剩余可读字节上限，初始值为“多次读取的总计限制数据量”
}

func (l *LimitedReader) Read(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, EOF
	}
	if int64(len(p)) > l.N { //!!! 如果缓存长度超过了剩余可读的上限，
		p = p[0:l.N] //!!! 则取长度为对应剩余可读上限的缓存子切片来作为读取操作的缓存
	}
	n, err = l.R.Read(p)
	l.N -= int64(n) //!!! 读取后消减剩余可读字节上限。
	return
}

// NewSectionReader returns a SectionReader that reads from r
// starting at offset off and stops with EOF after n bytes.

// NewSectionReader 返回一个SectionReader实例，SectionReader基于给定的“分段读取器”
// 对“特定分段”进行读取的读取器，程序员通过指定偏移（off）和读取字节数量（n），来设定“特定分段”。
// !!!虽然底层流还有数据可读，但“分段读取器”读取完所设定的分段后，就会以EOF结束。
func NewSectionReader(r ReaderAt, off int64, n int64) *SectionReader {
	var remaining int64 // 构建分段时，底层流r中剩余可读字节的位置,作为本段结尾位置。
	const maxint64 = 1<<63 - 1
	//如果 off + n  小于或等于最大的整数，则表明，流r中的off + n的位置是合法的。
	//remaining表示所读取流r中的最后一个字节的位置。
	if off <= maxint64-n {
		remaining = n + off
	} else {
		// Overflow, with no way to return error.
		// Assume we can read up to an offset of 1<<63 - 1.
		//如果 off+n 超过了最大的整数，那么，就将流r中的最后一个字节的位置设置为最大的整数，以防止超界访问。
		remaining = maxint64
	}
	return &SectionReader{r, off, off, remaining}
}

// SectionReader implements Read, Seek, and ReadAt on a section
// of an underlying ReaderAt.
// SectionReader是对给定的分段进行读取的读取器，通过指定偏移（off）和读取字节数量（n），来设定分段。
// !!!虽然底层流还有数据可读，但“分段读取器”读取完所设定的分段后，就会以EOF结束。
// SectionReader基于ReaderAt接口，在给定的“分段”实现了Read，Seek，和ReadAt方法。
// !!! SectionReader是对“分段读取器”的封装，使得“分段”读取更加方便。
type SectionReader struct {
	r    ReaderAt //具有分段读取能力的字节读取器，这里作为分段读取器的底层流。
	base int64    //分段读取器在底层流r中基址，也就是初始偏移量，永远不变。
	off  int64    //分段读取器在底层流r中的最新偏移量，也就是下一次读取的开始位置，每次读取后增加。
	//!!! 虽然构建分段流时指定了分段的偏移量off 与分段长度n，但是由于off+n可能超过最大整数，
	limit int64 //考虑到整数最大取值范围,结合偏移位置和分段长度，最终确定的分段读取器在底层流r中最大可读取字节的位置。
}

// 对Reader的实现
// !!!  分段读取器的Read方法可以多次调用，但是达到所在分段的结尾就会返回EOF。
func (s *SectionReader) Read(p []byte) (n int, err error) {
	//如果读取的偏移量大于最大可读取的字节位置，那么久表示已经读取完毕，没有字节可读了。
	if s.off >= s.limit {
		return 0, EOF
	}
	// 先计算最大可读字节数（max），为s.limit - s.off
	// 如果存放读取结果的内存缓存长度超过最大可读取的字节数 ,则截取长度正y8好满足
	// 最大的可读字节数的子切片来作为存储“分段读取器”r的内存缓存；否则，
	// 直接将长度不超过最大可读取字节数的给定缓存P来作为“分段读取器”r的内存缓存。
	if max := s.limit - s.off; int64(len(p)) > max {
		p = p[0:max] //截取刚好容纳所有剩余字节长度的子切片作为底层流r的读取缓存
	}
	//从off处读取len(p)个字节到缓存中p。
	n, err = s.r.ReadAt(p, s.off)
	//更新在底层流中的最新的偏移量。
	s.off += int64(n)
	return
}

var errWhence = errors.New("Seek: invalid whence")
var errOffset = errors.New("Seek: invalid offset")

// 是对Seeker接口的实现，意味着分段读取器可以跳跃、回退和重置。
func (s *SectionReader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	default:
		return 0, errWhence
	case SeekStart: //从头开始的偏移处读取
		offset += s.base
	case SeekCurrent: //从当前位置开始的偏移处读取。
		offset += s.off
	case SeekEnd: //!!! 从尾部开始的偏移处读取。 此时，有效的offset参数应为负值。
		offset += s.limit
	}
	if offset < s.base {
		return 0, errOffset
	}
	s.off = offset
	return offset - s.base, nil
}

// 是分段读取器对ReadAt（分段读取接口）的实现。
// 注意，输入的偏移量代表的是相对于分段读取器起始位置的偏移量。
func (s *SectionReader) ReadAt(p []byte, off int64) (n int, err error) {
	//如果给定的偏移超过了分段读取器对底层流的有效读取边界[s.base,s.limit)，则返回(0,EOF)
	if off < 0 || off >= s.limit-s.base {
		return 0, EOF
	}
	off += s.base //注意，现在的off代表的是相对于底层流的偏移位置。
	//如果buf的长度超过了最大可能读取的字节数量，也就是区间[s.base+off,s.limt)的长度，则截取
	//长度恰好为最大可能读取字节数的子切片来作为存储所读取数据的内存缓存。
	if max := s.limit - off; int64(len(p)) > max {
		p = p[0:max]
		n, err = s.r.ReadAt(p, off)
		if err == nil {
			err = EOF
		}
		return n, err
	}
	//从底层流的偏移位置处读取长度为len(p)个字节的数据。
	return s.r.ReadAt(p, off)
}

// Size returns the size of the section in bytes.
// Size方法返回分段的字节数量大小。
// !!! 虽然构建分段流时指定了分段的偏移量off 与分段长度n，但是由于off+n可能超过最大整数，
// !!! 所以，此时分段的长度为最大整数长度减去偏移量off，小于n。
func (s *SectionReader) Size() int64 { return s.limit - s.base }

// TeeReader returns a Reader that writes to w what it reads from r.
// All reads from r performed through it are matched with
// corresponding writes to w. There is no internal buffering -
// the write must complete before the read completes.
// Any error encountered while writing is reported as a read error.

// TeeReader返回一个字节读取器，该“字节读取器”的Read方法把从”底层字节读取器“r中所读取
// 的字节数据读入内存缓存的同时也”写入器“w。
// !!! 通过TeeReader所执行的所有来自r的读操作（Read方法的调用）必须与相应的（在内部）对w的写操作相匹配。
// !!! TeeReader的Read方法实现中没有中间的缓存——所有（读操作内部的）写操作一定会在读操作（Read方法调用）结束之前完成。
// !!! “读操作”内部向w写入数据的过程中所遇到的任何错误都会被当做“读操作（Read方法调用）”的错误进行报告。
// !!! 顾名思义，Tee是T型管，或者“三通”的意思，一次读取，同时把数据送入两个地方，即，
// !!!  内存和写入器。比如，使用TeeReader把从文件读取数据到内存的同时，也将其写入
// !!!  到标准输出设备进行显示就比较方便。
func TeeReader(r Reader, w Writer) Reader {
	return &teeReader{r, w}
}

type teeReader struct {
	r Reader
	w Writer
}

// !!!  对Reader接口的Read方法的实现，该Read方法实现了一次读取底层流，数据放入两个地方，
// !!!  内存缓存p，以及另一个流（通过w写入）。在Read方法结束前，两个地方的数据存放一定已经完成。
func (t *teeReader) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p) //将数据读取到内存缓存p中
	if n > 0 {           //!!!如果读到了数据，则会写入到w中，所以写入完成后，整个读操作才会完成。
		if n, err := t.w.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return
}

// Discard is a Writer on which all Write calls succeed
// without doing anything.
// Discard是一个写入器，所有的Write调用都会成功，而且什么都不做。
var Discard Writer = discard{}

type discard struct{}

// discard implements ReaderFrom as an optimization so Copy to
// io.Discard can avoid doing unnecessary work.
// 出于优化的目的，discard实现了ReaderFrom方法，这样Copy数据到
// io.Discard就可以避免面一些不必要的工作。
// !!! _是个可以抛弃的匿名变流量名， ReaderFrom是个接口类型，
// !!! 这样定义一个可抛弃的变量是想让编译器帮助程序员进行类型检查，也就是让编译器
// !!! 检查discard{}是ReaderFrom接口的实现。
var _ ReaderFrom = discard{}

// 对Wirter接口中的Write方法的实现 ，什么都不做，每次调用都成功。
func (discard) Write(p []byte) (int, error) {
	return len(p), nil
}

// 对StringWriter接口中的WriteString方法的实现 ，什么都不做，每次调用都成功。
func (discard) WriteString(s string) (int, error) {
	return len(s), nil
}

// !!!  定义一个全局共享的持有"黑洞"的缓存池，这里的黑洞是一个可以无限写入字节的切片，
// !!! 所谓无限写入字节，就是这个“定长切片”的中数据可以被反复地覆盖。
// !!! 之所以使用sync.Pool来缓存这个黑洞是因为sync.Pool可以保证多线程下的get操作
var blackHolePool = sync.Pool{
	New: func() any {
		b := make([]byte, 8192) // !!!“黑洞”是一个8k大小的切片,8k是UFS文件系统逻辑块（数据页）的默认大小。
		return &b
	},
}

// !!! 如果一个Writer实现ReadFrom方法，这样当使用Copy方法将Reader r中的数据
// !!! 写入到这个Writer中时，就会直接调用该方法，而不用在创建中间内存缓存来多次读写，
// !!! 详见Copy方法的实现。
func (discard) ReadFrom(r Reader) (n int64, err error) {
	bufp := blackHolePool.Get().(*[]byte) //从黑洞缓存池中获取存储数据的"黑洞"
	readSize := 0                         //重置Reader r每次读取操作所读取的字节数。
	for {                                 //  只要没有遇到err或者EOF，就不停地读取数据到“黑洞中”，并记录读取的总计字节数。
		readSize, err = r.Read(*bufp) //通过Reader r多次读取数据到“黑洞”这个切片中。
		n += int64(readSize)          //记录总计读取的字节数
		if err != nil {
			blackHolePool.Put(bufp)
			if err == EOF { //如果err是读到了流的结尾，
				return n, nil //就把err修正为nil。
			}
			return //返回命名的返回参数 （n,err）
		}
	}
}

// NopCloser returns a ReadCloser with a no-op Close method wrapping
// the provided Reader r.
// If r implements WriterTo, the returned ReadCloser will implement WriterTo
// by forwarding calls to r.
// !!! NopCloser在封装了给定的Reader r的基础，返回一个Close方法不做任何操作的ReadCloser的实现实例。
// !!! 如果给定的Reader r的具体实例还实现了 WriterTo接口（Writer调用客户端），那么返回的具体实例
// !!! 的类型为 nopCloserWriterTo，否则，返回具体实例的类型为nopCloser。
func NopCloser(r Reader) ReadCloser {
	if _, ok := r.(WriterTo); ok {
		return nopCloserWriterTo{r}
	}
	return nopCloser{r}
}

type nopCloser struct {
	Reader //!!! 嵌入了Reader，相当于继承，也就是该struct具有了Reader的功能。
}

func (nopCloser) Close() error { return nil }

type nopCloserWriterTo struct {
	Reader //!!! 嵌入了Reader，相当于继承，也就是该struct具有了Reader的功能。
}

func (nopCloserWriterTo) Close() error { return nil }

func (c nopCloserWriterTo) WriteTo(w Writer) (n int64, err error) {
	return c.Reader.(WriterTo).WriteTo(w)
}

// ReadAll reads from r until an error or EOF and returns the data it read.
// A successful call returns err == nil, not err == EOF. Because ReadAll is
// defined to read from src until EOF, it does not treat an EOF from Read
// as an error to be reported.

// ReadAll试图从“字节读取器r”中读取所有数据，直到遇到错误或者EOF，然后，返回它所读取的数据。
// 成功的调用会返回err==nil，而非err==EOF.因为，ReadAll被定义为从数据园中读取数据直到EOF，
// 所以，它不会把EOF当作错误来报告。
// !!! ReadAll 方法试图返回读取器中所有的有效数据。
func ReadAll(r Reader) ([]byte, error) {
	//!!! 512是物理块的大小，也就是（磁盘）控制器每次读取（一个扇区sector）数据的大小。
	b := make([]byte, 0, 512) //创建了一个len=0 ，cap=512的切片。
	for {                     //
		if len(b) == cap(b) { //如果当b的长度增长到与容量相同时，就扩容，扩容会产生一个新的底层内存。
			// Add more capacity (let append pick how much).
			// !!! 通过追加一个“0值”的字节来使存储结果的切片扩容，而扩容len+1，容量由append操作优化 ，
			// !!! ,扩容后使用取得子切片[:len(b)]操作使得b的长度仍然保持原来的长度，这样，后续的读取
			// !!! 追加操作就能够从扩容的第一个字节开始。
			b = append(b, 0)[:len(b)] // [:len(b)]中的b的len还是未扩容前的b的长度。
		}
		//将剩余容量部分的子切片作为Read操作的存储底层流数据的内存buffer。
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n] //根据读取的字节数修正b的边界为所读取的所有数据。
		if err != nil {
			if err == EOF {
				err = nil
			}
			return b, err
		}
	}
}
