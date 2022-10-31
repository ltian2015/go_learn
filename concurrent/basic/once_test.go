package basic

import (
	"sync"
	"testing"
	"time"
)

/*
*

sync包通过使用Once类型为在多个gotouintes存在的情况下进行性初始化提供了一个安全的机制
对于特定的 f,多个线程都可以执行once.Do(f)，但只有一个线程会运行f()，而其他调用
once.Do(f)的线程将会阻塞，直到到f()返回。
来自once.Do(f) 对f()的唯一调用，同步 先于任何所有对once.Do(f)
因为当多个线程同时调用时，只有一个线程能够执行f(),其他线程都会阻塞，等待f()的结束，
执行f()的线程也是要等待f()执行完毕才能继续 向下执行，所以“来自once.Do(f) 对f()的唯一调用，
同步 先于任何所有对once.Do(f)”，由于once.Do(f)要求f没有参数，也没有返回值，
所以，没有输入与输出的函数，可以看作一个独立的任务，在这个独立的任务一定是为了修改
独立于该任务之外的某些状态，否则，这样的任务就没有任何存在的意义。所以，这样的任务往往是一个
大工作的组成部分，大的工作定义了任务需要处理的外部状态。但是，这样的任务尽量少用，因为多个任务都
修改同一组外部状态，逻辑上并不好理解。所以，最好只用于初始化或者销毁工作。
*/
var wg sync.WaitGroup
var greetingWord string
var once sync.Once

func setup() {
	greetingWord = "hello ervybody"
	//如果setup() 被多次调用，那么应该打印多次，但这里，setup() 函数只在sync.Once.Do()中调用。
	println("setup task will over ,wait for 2 seconds")
	time.Sleep(time.Second * 2)
}
func doSetupA() {
	once.Do(setup) //这是一个潜在的等待点，只有一个线程可以执行，其他线程必须等待setup完成，然后越过等待点直接执行后面的代码
	println("A: ", greetingWord)
	wg.Done()
}
func doSetupB() {
	once.Do(setup)
	println("B: ", greetingWord)
	wg.Done()
}
func doSetupC() {
	once.Do(setup)
	println("C: ", greetingWord)
	wg.Done()
}

func TestOnce(t *testing.T) {
	println("begin test")
	wg.Add(3)
	go doSetupA()
	go doSetupB()
	go doSetupC()
	wg.Wait()
}
