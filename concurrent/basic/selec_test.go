package basic

import (
	"fmt"
	"testing"
	"time"
)

// 学习select用法
/**
!!! select 语句可以让Go 例程“阻塞等待”多个“通道操作（send/receive）”的完成。

语法上，select 语句用于选择从多个 send/receive"通道操作"中做出选择。每个通道操作都是一个case表达式，
一个case表达式对应一个case代码段。case表达式与case代码段之间用冒号（:）分割，语法是：
case [case 通道send/receive表达式] :
	 case代码段开始
	   .......
	 case 代码段结束
case [case 道send/receive表达式]] :
	  case代码段开始
	   .......
	 case 代码段结束
......

!!!select语句在程序运行时会自上至下对每个“通道操作（send/recieve）的case表达式”进行求值，
如果发现只有一个case表达式中的“通道操作send/recieve”已经就绪（通道不阻塞），就执行该分支。
如果发现有多个分支中的“通道通信操作”都已经完成（就绪）时，就会随机选择一个分支执行。


因此，Select语句会产生阻塞，一直阻塞到 “其中有一个” send/receive操作已经成功，就执行该就绪的通道操作下case代码段。
如果同时多个操作都已经成功，那就随机从中选择一个通道操作。
!!!如果存在default，当没有通道操作成功的时候，就会直接执行default分之,这样，select语句所在例程不会阻塞。
!!!虽然，default段的存在能够避免死锁（deadlock）的发生，但也会使得阻塞分支的大概率不会得到执行，因为select语句执行时，
!!!很可能"通信"
除了每个case分支的表达式都是一个通道操作（channel operation）之外，select语法于switch 语法非常相似。
虽然不是每个case都被执行，但是会从上到下对每个case表达式求值。
*/
func LearnSelect1() {

	ch1 := make(chan interface{})
	ch2 := make(chan int)
	ch3 := make(chan int)
	//直接在goroutine中运行匿名的函数字面量
	go func() {

		time.Sleep(1 * time.Second)
		ch1 <- 1
	}()

	go func() {

		time.Sleep(4 * time.Second)
		ch2 <- 4
	}()

	go func() {
		time.Sleep(1 * time.Second)
		ch3 <- 18
	}()

	fmt.Println("Blocking on read..., if there's no default ")
	//如果有default段情况，则所有case段失效。
	select {
	case <-ch1: //case1 表达式，表达一种通道操作情况
		//case1 代码段
		fmt.Println("ch1 case...")

	case v := <-ch2: // case1 表达式，表达一种通道操作情况
		//case1 代码段，
		fmt.Println("ch2 case...")
		fmt.Println(v)
	case v := <-ch3: // case3 表达式，表达一种通道操作情况
		//case3 代码段
		fmt.Println("ch3 case...")
		fmt.Println(v)
		//有了defualt语句，则上述case分支大概率得不到执行，只执行default代码段。
		// 但是，每个case 表达式都会被求值。详细情况可以看 func LearnSelect2()
	default:
		fmt.Println("ch2 case...")
		fmt.Println(<-ch2)
	}
}

var ch1 chan int
var ch2 chan int
var chs = []chan int{ch1, ch2}
var numbers = []int{1, 2, 3, 4, 5}

func SelectCaseEvaluate() {
	//这个例子表明，虽然不是每个case代码段都会被会执行，但是每个case 表达式都会被求值（执行）。
	select {
	case getChan(1) <- getNumber(1): //case 1表达式
		//case 1 代码段
		fmt.Println("1th case is selected.")
	case getChan(2) <- getNumber(2): //case 2表达式
		//case 2 代码段
		fmt.Println("2th case is selected.")
	default:
		//default 代码段
		fmt.Println("default is select!.")
	}
}

func getNumber(i int) int {
	fmt.Printf("got numbers[%d]\n", i)

	return numbers[i]
}
func getChan(i int) chan int {
	fmt.Printf("got chs[%d]\n", i)
	return chs[i-1]
}

// 空的select 语句因为没有通道操作，也会导致死锁异常发生（panic）
func DeadlockSinceNoChannelOp() {
	println("程序会发生死锁,因为空select语句没有通道操作，所以会一直等待，导致所有例程休眠！")
	select {}
}

// 当通道操作永远不会关闭的时候，select就会发生死锁异常（panic）
func TestDeadlockSinceChannelOpNeverOccured(t *testing.T) {
	var ch chan string = make(chan string)
	println("程序会发生死锁,因为在等待通道的操作的发生，而没有任何例程对等待的通道进行操作！！")
	select {
	case <-ch:
	}
}

// !!!通过select实现超时控制
// 通过time.After函数返回一个只读的channel，向里面发送一个当前的时间，
// 这样就可以实现超时的功能。
// 只要在select语句中使用time.After函数返回的channel，就可以实现超时。
// 这个函数返回的channel是只读的，不能向里面发送数据，只能接收数据。
// 这个函数的返回值是一个只读的channel，里面发送的是当前的时间。
func TestImplementTimeoutWithSelect(t *testing.T) {
	chWork := make(chan int)

	select {
	case v := <-chWork:
		fmt.Println("work done, value:", v)
		//!!! time.After 函数返回了一个只读的channel，并向里面发送了一个当前的时间。
	case timeOnTimeout := <-time.After(10 * time.Second):
		fmt.Println("timeout 10 seconds!, current time is :", timeOnTimeout)
	}
}
