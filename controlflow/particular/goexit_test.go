package particular

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

/**
    Goeixt与panic都使得函数非正常退出，而且都会向外层函数传播不同的信号。
	Goexit使函数与一个Goexit信号关联，而且无法像panic所发出的panic信号那样
	 可以通过recover机制来消除。因此会使得Goeixts一定会导致goroutine入口函数退出，
	 所以，Goexit是goroutine退出命令。
	 且由于二者都是一种非正常的退出机制，在一个Goroutine中，
	 函数正常执行阶段Goexit与Panic两种退出机制不能连续发生，也不能同时发生。并且，
	 Goexit调用会导致函数进入退出阶段，从而也引发被推迟执行的函数调用。

	 由于在正常执行阶段的Goexit不是panic，也不能与之并存，所以当GoExit引发函数退出时，
	 在任何一个被推迟的函数中调用recover都只能得到nil。

	 但通过defer机制，可以使得panic与Goexit一个在正常阶段执行，一个在退出阶段执行，
	 就会使得二者相继发生。如果不了解其机制就会发生bug。
	 如果先执行panic（正常执行阶段），后执行Goexit（defer到退出阶段执行），后执行的Goexit会掩盖panic。
	 如果，先执行Goexit（正常执行阶段），再执行panic（defer到退出阶段执行），先执行的panic就不会被Goexit掩盖了。

	 在一个普通的goroutine中调用Goexit函数会结束该goroutine，但其他goroutine不会受影响.
	 主goroutine的入口函数（main函数）会对“go运行时”返回状态码，而普通goroutine的入口函数没有返回值。
	 这是主goroutine的入口函数main与普通goroutine入口函数的区别。

	 但是，如果在主 goroutine（main函数为入口的goroutine）中使用了Goexit，会导致了主线程的提前结束，
	 由于入口main函数没有返回状态结果，因此，在“go运行时”看来，主routine处于崩溃状态。

	 但要注意的是，在调用Goexit结束主goroutine时，其他正在执行子goroutine
	 不会受影响，仍然会继续执行，直至终结。
**/
func TestGoexitInMainGoroutine(t *testing.T) {
	//程序入口函数中推迟执行的匿名函数
	defer func() {
		fmt.Println("function will terminate!") //会打印。
	}()
	runtime.Goexit() //测试运行环境中主goroutine，会导致测试运行崩溃。
	fmt.Println("unreachable")
}

//测试使用goexit结束子goroutine，
//这里有三个子goroutine，前两个很快就会执行，并通过调用Goexit退出，第三个
//则先休眠3秒，然后再打印一句话后正常退出。可见，Goexit调用使得各自子goroutine退出
//不会对其他执行中的goroutine产生影响。
func TestGoexitInSubGoroutine(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		fmt.Println("execute task1 in goroutine :", getGoroutineID())
		wg.Done()
		runtime.Goexit() //普通的goroutine使用Goexit()结束不会产生崩溃。

	}()
	go func() {
		fmt.Println("execute task2 in goroutine :", getGoroutineID())
		wg.Done()
		runtime.Goexit() //普通的goroutine使用Goexit()结束不会产生崩溃。

	}()
	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("execute taks 3 in goroutine :", getGoroutineID())
		wg.Done()

	}()
	wg.Wait()
	fmt.Println("main programe exit ,it 's run in goroutine: ", getGoroutineID())
}

//获得GOroutine ID,调试时使用，不应用在正式代码中。
func getGoroutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

//TestGoexitTogetherWithPanic函数测试了Panic与Goexit同时在一个函数中使用的情况，
//正常执行阶段，两个并行的机制应该互不干扰，也不会相互掩盖。但是事实是，在退出阶段，
//如果先执行panic，后执行Goexit，后执行的Goexit会掩盖panic。
//如果，先执行Goexit，再执行panic，panic就不会被Goexit掩盖了。

//case1代码中，先执行的panic被后执行的Goexit掩盖了，所以不会被recover捕获.
//case2代码中，先执行的Goexit没有影响后执行的panic，panic可以被recover捕获到。
//case1如下：
func TestGoexitTogetherWithPanicCase1(t *testing.T) {
	go func() {
		// 先抛出的panic被后执行的Goexit所掩盖
		// recover()将不会捕获到panic
		defer func() {
			r := recover() //r==nil
			if r != nil {
				fmt.Println("catch the panic : ", r)
			}
		}()
		defer runtime.Goexit() //后执行Goexit
		panic("bye")           //先抛出panic
	}()
	//下面的语句将会等待系统中存活线程数变为1才会结束循环。
	//主线程会一直以非阻塞的方式等待另一个线程的完成。
	//这是一种非阻塞式的线程等待方式，它不会一直占用CPU资源的自旋锁(spin lock)。
	//Gosched()与for循环配合是实现自旋锁的重要机制，自旋条件可以根据情况设定。
	//这里是存活的线程数量，自旋锁更多场景是CAS（Compare And Swap）算法中。
	for runtime.NumGoroutine() > 2 {
		fmt.Println("there are ", runtime.NumGoroutine(), " goroutines")
		//当前goroutine让出所获得的CPU执行权，让其他goroutine执行，
		//但是,会自动重新获取到CPU的执行权,获得执行权后，再从此语句之后的一个语句开始执行。
		//由于在for循环中，再次获得执行权后，会从for循环的循环条件判断语句开始执行，
		//也就是 判断 runtime.NumGoroutine() > 2是否成立，如果成立就进入再次循环，不成立
		//就退出循环。
		//注意，测试环境中，会多出一个goroutine，故而用runtime.NumGoroutine() > 2
		//在使用main()函数的正式运行环境中，不会多出一个goroutine，故而应该用
		//runtime.NumGoroutine() > 1的自旋条件。
		runtime.Gosched() //相当于Java中Thread对象的yield方法
	}
	fmt.Println("main goroutine will exit!") //会被执行
}

//case2如下：
func TestGoexitTogetherWithPanicCase2(t *testing.T) {
	go func() {
		defer func() {
			r := recover() //r!=nil,可以被捕获到
			if r != nil {
				fmt.Println("catch the panic : ", r)
			}
		}()
		defer panic("bye") //后抛出panic
		runtime.Goexit()   //先执行Goexit()
	}()
	//下面的语句将会等待系统中存活线程数变为1才会结束循环。
	//主线程会一直以非阻塞的方式等待另一个线程的完成。
	//这是一种非阻塞式的线程等待方式，它不会一直占用CPU资源的自旋锁(spin lock)。
	//Gosched()与for循环配合是实现自旋锁的重要机制，自旋条件可以根据情况设定。
	//这里是存活的线程数量，自旋锁更多场景是CAS（Compare And Swap）算法中。
	for runtime.NumGoroutine() > 2 {
		fmt.Println("there are ", runtime.NumGoroutine(), " goroutines")
		//当前goroutine让出所获得的CPU执行权，让其他goroutine执行，
		//但是,会自动重新获取到CPU的执行权,获得执行权后，再从此语句之后的一个语句开始执行。
		//由于在for循环中，再次获得执行权后，会从for循环的循环条件判断语句开始执行，
		//也就是 判断 runtime.NumGoroutine() > 2是否成立，如果成立就进入再次循环，不成立
		//就退出循环。
		//注意，测试环境中，会多出一个goroutine，故而用runtime.NumGoroutine() > 2
		//在使用main()函数的正式运行环境中，不会多出一个goroutine，故而应该用
		//runtime.NumGoroutine() > 1的自旋条件。
		runtime.Gosched() //相当于Java中Thread对象的yield方法
	}
	fmt.Println("main goroutine will exit!") //会被执行
}
