package basic

import (
	"sync"
	"testing"
)

/*
*
once.Do(f)函数当且仅当 Do 是首次针对“此 [Once]” 实例调用时，Do 才会调用函数 f。换句话说，给定
var once Once
如果 once.Do(f) 被多次调用，即使 f 在每次调用中都有不同的值，也只有第一次调用会调用 f。
因此，如果想要函数每次调用都返回不同的值，那么每次执行函数时都需要一个新的 Once 实例。
Do 用于初始化，并且必须只运行一次。由于 f 是 niladic（空参数）类型的函数，
因此可能需要使用“函数字面量”来捕获 Do 调用函数的参数：
config.once.Do(func() { config.init(filename) })

在 f 的一次调用返回之前，所有对Do的（并发）调用都不会返回，因此，
如果 f 在执行中又导致了同一个【Once】实例的Do函数被调用，就会导致将死锁，
这两个情况都是由于Do在调用f时使用了锁，只有第一个获得锁的调用才会执行f，其他的调用都必须等待。
在f中又调用Once】实例的Do函数就会导致这个执行线程又来请求被自己获得锁，这就是死锁。

如果 f 发生panic，Do 会认为它已返回；以后对 Do 的调用将返回，而不会调用 f。
这是因为f中的panic会导致f进行入f生命周期的退出阶段，所以Do认为它已经返回。

sync包通过使用Once类型为在多个gotouintes存在的情况下进行性初始化提供了一个安全的机制
对于特定的 f,多个线程都可以执行once.Do(f)，但只有一个线程会运行f()，而其他调用
once.Do(f)的线程将会阻塞，直到到f()返回。
来自once.Do(f) 对f()的唯一调用，同步 先于任何所有对once.Do(f)
因为当多个线程同时调用时，只有一个线程能够执行f(),其他线程都会阻塞，等待f()的结束，
执行f()的线程也是要等待f()执行完毕才能继续 向下执行，所以“来自once.Do(f) 对f()的唯一调用，
同步 先于任何所有对once.Do(f)”，由于once.Do(f)要求f没有参数，也没有返回值，必须通过
函数字面量来捕获参数并回写结果。

*/

var greetingWord string
var once sync.Once  //!!! 同一个函数被同一个sync.Once实例多次调用，只会实际执行一次，
var once2 sync.Once //!!! 同一个函数在不同sync.Once实例中调用，每个Once实例会执行一次。

func initialize() {
	greetingWord = "hello ervybody"
	//如果initialize() 被几次调用，那么下面的这段话就会被打印几次。
	println("初始化工作已经完成！")
}
func doInitializeA() {
	once.Do(initialize) //!!! 这是一个隐式的阻塞点，只有一个线程可以执行，其他线程必须等待initialize完成，然后越过等待点直接执行后面的代码
	println("A: ", greetingWord)

}
func doInitializeB() {
	once.Do(initialize)
	println("B: ", greetingWord)

}
func doInitializeC() {
	once2.Do(initialize)
	println("C: ", greetingWord)

}

func TestOnce(t *testing.T) {
	var wg sync.WaitGroup
	println("begin test")
	wg.Go(doInitializeA)
	wg.Go(doInitializeB)
	wg.Go(doInitializeC)
	wg.Wait()
}
