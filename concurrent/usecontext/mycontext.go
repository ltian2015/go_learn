package usecontext

import (
	"time"
)

// A Context carries a deadline, a cancellation signal, and other values across
// API boundaries.
// Context 携带截止时间（deadline），取消信号（cancellation signal） 以及其他穿过API边界的值。
//
// Context's methods may be called by multiple goroutines simultaneously.
// Context的方法可以被多个例程（goroutines）同时调用。
type MyContext interface {
	// Deadline returns the time when work done on behalf of this context
	// should be canceled. Deadline returns ok==false when no deadline is
	// set. Successive calls to Deadline return the same results.
	//Deadlinne函数根据该上下文进行的工作何时应当取消的时间。当没有设置deadline时间时，
	//Deadline函数返回ok==false。后续对Deadline方法的调用都会返回相同的结果。
	Deadline() (deadline time.Time, ok bool)

	// Done returns a channel that's closed when work done on behalf of this
	// context should be canceled. Done may return nil if this context can
	// never be canceled. Successive calls to Done return the same value.
	// The close of the Done channel may happen asynchronously,
	// after the cancel function returns.
	// DoneHanns函数返回一个通道（channel),当根据该上下文进行的工作被取消时，该通道是关闭的。
	// 如果上下文还不能被取消，那么Done方法就会返回nil。后续的方法调用都会返回相同的值。
	// 在cancel函数返回后，会异步产生Done通道的关闭操作。
	//
	// WithCancel arranges for Done to be closed when cancel is called;
	// WithDeadline arranges for Done to be closed when the deadline
	// expires; WithTimeout arranges for Done to be closed when the timeout
	// elapses.
	// WithCancel方法就是为了调用cancel关闭Done通道而准备的；
	// WithDeadline 就是为了当超过截止日期关闭Done通道而准备的；
	// WithTimeout 就是为了超时是关闭Done通道而准备的;
	// Done is provided for use in select statements:
	//
	//  // Stream generates values with DoSomething and sends them to out
	//  // until DoSomething returns an error or ctx.Done is closed.
	// Done函数主要用于select语句
	//  // 下面的Stream函数使用DoSomething并将产生一些值，并在DoSomething返回错误或者
	//  // ctx.Done被关闭时，将这些值法发送到外面。
	//  func Stream(ctx context.Context, out chan<- Value) error {
	//  	for {
	//  		v, err := DoSomething(ctx)
	//  		if err != nil {
	//  			return err
	//  		}
	//  		select {
	//  		case <-ctx.Done():
	//  			return ctx.Err()
	//  		case out <- v:
	//  		}
	//  	}
	//  }
	//
	// See https://blog.golang.org/pipelines for more examples of how to use
	// a Done channel for cancellation.
	//见https://blog.golang.org/pipelines 以获得更多关于如何使用Done 通道用于取消操作的的例子。
	Done() <-chan struct{}

}


