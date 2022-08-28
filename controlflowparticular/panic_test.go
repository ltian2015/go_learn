/**
defer和panic/recover是go语言中特殊的程序控制流。
这里主要关注GO语言的error机制、panic异常抛出机制、defer（延迟执行）机制，
以及defer与recover结合捕获并消除panic异常，实现程序从异常中恢复的机制。
本文件重点关注panic机制，及panic的恢复机制。尤其是多线程（goroutine）下
的panic影响范围与恢复机制。
**/
package controlflowparticular

/**
1. 如何认识GO语言中的error与panic的区别。
GO语言函数提供了两种报错机制，一种是返回error机制，一种是抛出panic机制。
返回error，表示函数执行中出现的error属于程序正常的领域业务处理逻辑考虑范围之内的正常业务问题，
将其以error数据的方式返回给外层调用程序，有助于外层程序对error做出针对性的领域业务逻辑处理，
这些error属于程序正常领域业务逻辑范围要考虑的问题。比如，取款金额超过了账户余额，
应该返回error而不是panic。

而抛出panic类似于Java的异常机制，表明函数中出现的问题与程序所在领域的业务处理逻辑无关，往往属于bug，
或系统级的问题，由于这类问题不在程序所在领域的业务处理逻辑考虑范围之内，所以不应以函数的调用结果的方式
，也就是error返回给外层函数处理。这类问题就是panic，对panic的处理往往是只能记录问题，交给系统监控者或者开发者修改bug。
比如，数组访问超界、网络连接故障、对不允许拷贝的对象进行了拷贝、预想的运行环境没有正确初始化等。

举一个现实场景的例子，表明error与panic的应用。
比如，对于一个银行账户处理领域中，有一个函数可以根据给定的账户号参数和钱数，从账户中取出给定数额的钱，函数返回余额。
在这个问题域中，考虑到取钱余额不足是正常的逻辑。所以，如果该函数的调用着所给定的账户号所指定的账户余额不足，
则应返回一个error，这样调用者可以对该类型的错误进行相应处理，比如换一个账户进行操作，或者生成给操作者的人性化提示信息。
但如果访问数据库失败，这就不是该领域所应处理的问题，则应抛出异常，因为这是来自另一个正交领域的问题。
也应在另一个正交领域解决panic，比如，尝试重新连接数据库，或者写入异常日志。

GO语言将领域内的业务逻辑错误（error）从广泛意义的错误中“显式地”剥离出来，有助于使业务逻辑处理更加清晰，
使得领域业务逻辑中对业务逻辑错误的正常处理，与领域之外的各种异常的处理分离开来。

2.GO语言中的panic的抛出与捕获机制。

在GO语言中，可以使用内置的panic函数抛出一个“异常”，异常被抛出时，调用panic的函数，假定是f1_1，立即终止正常执行，
进入退出阶段。如果在退出阶段没有使用内置的recover函数捕获并清除f1_1函数的panic状态，那么，执行完退出阶段后，
这个panic会继续抛出到所在的f1_1函数的外层函数（假定是f1)调用f1_1之处，从此处终止函数f1的正常执行，使得
函数f1进入到退出阶段，这称之为panic的“传播”。如果函数f1在退出阶段也没有捕获该panic，那么这个panic就会继续向
上层函数“传播”，导致上层函数也进入到退出阶段。如此，一直抛出到直到goroutine的"入口函数"，
如果goroutine的“入口函数”在退出阶段也没有捕获并处理panic，就会导致（可能存在多goroutine）程序的整体崩溃。
更具体地细节是：
当函数f调用panic函数发出时候一个panic p1的时候，函数f就会与panic p1关联在一起，
如果在退出阶段，函数f又收到了一个新的panic p2，函数f就会与p2关联在一起，被关联的p1就会被
p2替换掉。由于函数调用通常是嵌套调用的，如果内层函数调用返回时关联了一个panic（没有使用recover移除），
这个panic就与外层函数关联，成为外层函数的panic，这样一个函数无论嵌套调用多少层，每个函数最多只能
关联一个panic，也就只能捕获这个panic。又考虑到goroutine的执行实际上由入口函数所引发的嵌套式的函数调用，
因此一个gotouine中最多只有一个panic与之关联。

由此，关于panic的抛出与捕获机制有以下4点关键：
2.1 go语言内置的panic函数调用会导致所在函数立即结束正常代码执行而进入退出阶段。
2.2 只能通过defer机制，在函数退出阶段使用go语言内置函数recover函数捕获此前（自身或被嵌套调用函数）所
    抛出的panic。recover函数会移除goroutine中的panic状态，使得goroutine恢复正常状态，而不会崩溃。
2.3 panic(nil)也会导致goroutine处于panic状态，不使用recover()函数移除panic状态，goroutine仍会崩溃，
    但是，recover()函数所捕获的panic值是nil。
2.4 无法跨goroutine捕获panic，如果程序开启了多个goroutine，必须保证每个goroutine实例都能够捕获处理
    在自己执行过程中所抛出的panic，否则就可能因某个goroutine抛出panic而导致整个程序的崩溃。
	整体崩溃的后果是，其他正在执行的goroutine在整个程序崩溃时刻还没来得及执行的代码都不会被执行，
	因此后果是很严重，对多goroutine的程序而言，一个goroutine发生了panic，会殃及所有正在执行中的goroutine，
	导致所有goroutine的执行情况不可预料。所以，对于并发程序而言，每个gorouine的入口函数都是该goroutine
	实例捕获其内panic的最后一道防线，必须非常注意。
2.5 想要使用recover()捕获得到函数f（或其嵌套调用的函数）中发生的panic，
   recover()必须在正确的位置调用，也就是在f函数的一个defer函数,假定为panicHandleFn,代码中直接调用
   recover()才能生效。即：
  rule1. 下函数f中直接defer recover()函数无法消除函数f中的panic。
  rule2. 在被defer 的panicHandleFn函数的任何嵌套函数中调用recover，无法消除defer panicHandleFn
         语句所在函数f中的panic。
  rule3. 在函数f的非defer函数及其任何嵌套函数中调用recover函数，无法消除函数f中发生的panic。
  所以，函数f中，可以捕获并消除（传播到）f中的panic唯一有效的reover调用位置如下所示：
  func f(){
         defer func(){ //
			 r:=recover()  // 唯一正确位置
			 if r!=nil{

			 }
		 }
		 ....
		 panic("a panic")
  }
  2.6 还需要注意的是，对于标准的go编译器所编译的程序，有一些致命的问题发生后，可能无法用recover消除panic而恢复正常。这些
      致命问题包括堆栈溢出，内存超界。
**/

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

//展示单协程（goroutine）中抛出panic，导致函数异常中断执行，进入退出阶段的情景。
func TestPanicInOneGoroutine(t *testing.T) {
	defer func() {
		println("programe exit here") //会被打印。
	}() //将在由于panic所导致的退出阶段时执行该函数。

	println("hi, programe begin!")
	panic("bye") //导致所在函数进入退出阶段

	println("unreachable ") //这一行执行不到了。因为前面发生了panic
}

//展示单协程中抛出panic，并用defer+recover恢复异常的情景。
func TestPanicRecoverInOneGoroutine(t *testing.T) {
	defer func() {
		println("programe exit here") //会被打印。
	}() //将在由于panic所导致的退出阶段时执行该函数。
	defer func() {
		v := recover()
		fmt.Println("recovered:", v)
	}() //将在由于panic所导致的退出阶段时执行该函数,使得主函数从panic中恢复。
	println("hi, programe begin!")
	panic("bye") //导致所在函数进入退出阶段

	println("unreachable ") //这一行执行不到了。因为前面发生了panic
}

//测试多个goroutine中发生panic对其他正常goroutine和整个程序的影响。
//并发执行两个不相关的子任务，但是在任务较少，执行时间短的子任务2中抛出panic。
//会对执行任务更重，执行时间更长（用休眠3秒）子任务1产生影响。

func TestPanicAffectInMultiGoroutines(t *testing.T) {
	defer func() {
		r := recover()
		fmt.Println("recover panic : '", r, "' is catched in main goroutine!")
	}() //试图在主goroutine中捕获并恢复子goroutine中发生的panic，保持程序正常而不会崩溃，但这个努力是徒劳的！

	fmt.Println("main task begin")

	wg := sync.WaitGroup{}
	wg.Add(2) //为主任务增加两个并发执行的子任务。只有两个任务都完成，主任务才能继续工作。
	//子任务1，任务重，执行难时间长（用休眠3秒钟模拟）
	go func() { // 子routine的入口函数。子routine内没有被recover的panic最终会抛给此函数
		defer wg.Done()                   //等待计数减1。可能会被执行，取决于另一个goroutine的发生panic导致整体崩溃的发生时刻
		fmt.Println("sub task 1   begin") //可能会被执行，取决于另一个goroutine的发生panic导致整体崩溃的发生时刻
		func() {
			time.Sleep(3 * time.Second)    //在此休眠期间，由于另一个goroutine的panic所导致的整体崩溃，使得本goroutine后续代码无法接续执行。
			fmt.Println("sub task 1  end") //其他goroutine发生了panic，导致整体崩溃，本行代码来不及执行。
		}()
	}()
	//子任务2，任务较轻，执行时间短
	go func() { // 子routine的入口函数。子routine内没有被recocer的panic最终会抛给此函数。
		defer wg.Done()                     //等待计数减1,会被执行
		fmt.Println("sub task 2  executed") //会被执行
		func() {
			fmt.Println("sub task 2  will panic") //会被执行
			panic("sub-taks2-panic")              //抛出 panic，因为，没有被捕获，会导致程序整体崩溃，殃及另一个goroutine
		}()
	}()
	wg.Wait() //等待并发执行的所有子任务完成。
	//由于阻塞等待期间，子任务2的goroutine抛出了panic，导致整个程序崩溃,因此，下面代码永远不会被执行。
	fmt.Println("can not reachable")
}

//测试通过defer和recover机制保护每个goroutine实例不会因发生panic而崩溃整个程序。
// 为每个goroutine实例增加一个守候函数是一个很好的实践。
func TestPanicRecoverInMultiGoroutines(t *testing.T) {
	//定义一个goroutine保护函数，通过调用recover(）函数，保护goroutine不会因panic而引起整个程序的崩溃
	goroutineProtecter := func(taskName string) {
		r := recover()
		if r != nil {
			fmt.Println("recover panic: '", r, "' in ", taskName)
		}
	}

	fmt.Println("main task begin")
	wg := sync.WaitGroup{}
	wg.Add(2) //为主任务增加两个并发执行的子任务。只有两个任务都完成，主任务才能继续工作。
	//子任务1，任务重，执行难时间长（用休眠3秒钟模拟）
	go func() { // 子routine的入口函数。子routine内没有被recover的panic最终会抛给此函数
		defer wg.Done()                        //会被执行
		defer goroutineProtecter("sub task 1") //保护goroutine
		fmt.Println("sub task 1   begin")      //会被执行
		func() {
			time.Sleep(3 * time.Second)    //会被执行
			fmt.Println("sub task 1  end") //会被执行
		}()
	}()
	//子任务2，任务较轻，执行时间短
	go func() { // 子routine的入口函数。子routine内没有被recover的panic最终会抛给此函数。
		defer wg.Done()                        //会被执行
		defer goroutineProtecter("sub task 2") // 保护goroutine
		fmt.Println("sub task 2  executed")    //会被执行
		func() {
			fmt.Println("sub task 2  will panic") //会被执行
			panic("sub-task2-panic")              //抛出 panic，将会被goroutine的入口函数的在退出阶段捕获处理。
		}()
	}()
	wg.Wait()                     //等待并发执行的所有子任务完成。
	fmt.Println("will reachable") //尽管某个子goroutine中发生了panic，但是都被捕获处理了，不会影响整个程序
}

/**
   想要使用recover()捕获得到函数f（或其嵌套调用的函数）中发生的panic，
   recover()必须在正确的位置调用，也就是在f函数的一个defer函数（假定为deferedFn）的代码中直接调用
   recover()才能生效。即：
  rule1. f函数中不能直接defer recover函数
  rule2. 也不能在f函数被defer的函数（假定为deferedFn）的任何嵌套函数中调用。
  rule3. 不能在函数f的非defer函数及其任何嵌套函数中调用recover函数。
  函数f中，可以捕获并消除（传播到）f中的panic，在f中使用recover唯一有效的位置如下所示：

  func f(){
         defer func(){
			 r:=recover  // 唯一正确位置
			 if r!=nil{

			 }
		 }
		 ....
		 panic("a panic")
  }
**/
//给出一种典型的无效的Recover位置
func TestTypicalNoEffectRecvoerPostion(t *testing.T) {
	defer recover() //无效的recover 位置。
	panic(1)
}

//下面例子给出了大多数典型的无效defer位置
func TestMostNoEffectRecoverPositions(t *testing.T) {
	defer func() {
		defer func() {
			recover() // 无效，违反rule2。在defer函数的下级嵌套函数中调用，不是直接调用。
		}()
	}()
	defer func() {
		func() {
			recover() /// 无效，违反rule2。在defer函数的下级嵌套函数中调用，不是直接调用。
		}()
	}()
	func() {
		defer func() {
			recover() // 无效，违反rule3,rule2。在非f的defer函数中调用，而且不是直接调用，而是在下级嵌套函数中调用
		}()
	}()
	func() {
		defer recover() // 无效，违反rule3，rule1。 在非defer函数中，defer recover
	}()
	func() {
		recover() // 无效，违反rule3，在非defer函数中调用recover
	}()
	defer recover() // 无效，违反rule1，直接defer recover
	panic("bye")
}

//基于Recover的正确位置，下面演示了嵌套panic的recover顺序
func TestRecoverOrder(t *testing.T) {
	defer func() {
		fmt.Println("recover the panic ", recover()) //首先恢复panic 1
		defer func() {
			fmt.Println("recover the panic ", recover()) //然后恢复panic 2
		}()
		defer recover() // 无效，无法恢复任何panic
		panic(2)        //会被panic（2）所在函数的defer函数代码中的recover恢复。
	}()
	panic(1) //会被panic(1)所在函数的defer函数代码中的recover恢复。
}

//演示如果Panic抛出了nil 情况。
//程序中可以抛出nil panic，但如果不用recover捕获该panic，那么
//程序仍然会因为存在panic装态而崩溃，尽管此种情况下，recover捕获到得panic会是nil，
//它仍然阻止了程序的崩溃。
func TestPanicNil(t *testing.T) {
	defer func() {
		r := recover() //消除了goroutine的panic状态，使得goroutine不会崩溃
		if r != nil {  //由于panic是nil，所以，下面代码不会执行
			fmt.Println("recover from panic ", r)
		}
	}()
	panic(nil)             //尽管抛出的panic是nil，但是如果不被recover，程序仍会崩溃。
	println("unreachable") //虽然抛出了了异常值是nil，还是会到导致函数在panic处立即结束
}

/**
当函数f调用panic函数发出时候一个panic p1的时候，函数f就会与panic p1关联在一起，
如果此时（处于退出阶段）函数又收到了一个新的panic p2，函数f就会与p2关联在一起，被关联的p1就会被
p2替换掉。由于函数调用通常是嵌套调用的，如果内层函数调用返回时关联了一个panic（没有使用recover移除），
这个panic就与外层函数关联，成为外层函数的panic，这样一个函数无论嵌套调用多少层，最多只能
关联一个panic，也就只能捕获这个panic。
又考虑到goroutine的执行实际上由入口函数所引发的嵌套式的函数调用，因此一个gotouine中最多只有一个panic
与之关联。
**/
//TestPanicAssociationOrder函数测试panic的关联的替换顺序。
func TestPanicAssociationOrder(t *testing.T) {
	// 第四个被执行的延迟调用的函数，它所看到的是主函数当前关联的panic是3.
	defer func() {
		//recover函数移除并返回了主函数所关联的panic 3,这样主函数就处于正常状态。
		//这样，运行该主函数的goroutine就不会因为返回的主函数带有panic而崩溃。
		fmt.Println(recover()) //打印的是3。
	}()

	defer panic(3) // 第三个被执行的延迟调用函数，它所抛出的panic 3 替换主函数所关联的panic 2
	defer panic(2) // 第二个被执行的延迟调用函数，它所抛出的panic 2 替换主函数所关联的panic 1
	defer panic(1) // 第一个被执行的延迟调用函数，它所抛出的panic 1 替换主函数所关联的panic 0
	panic(0)       //在主函数正常执行阶段所抛出的panic  0，因此，成为首先被主函数关联panic 0
}

//演示goroutine中运行的函数与panic关联的顺序。
func TestPanicAssociationOrderInGoroutine(t *testing.T) {
	// 开启新的goroutine
	go func() {
		/**
		defer func() {
			fmt.Println(recover())   // 打印 2
		}() //捕获并移除 panic 2
		**/

		//一个被推迟调用匿名函数。当这个匿名函数完全退出的时候，panic 2 将会传播
		//给新goroutine的入口函数,并替换了入口函数原有关联的panic 0。
		//panic 2 没有被recover。
		defer func() {
			//panic 2将会替换你panic 1
			defer panic(2)
			//调用了一个匿名函数，当这个匿名函数完全退出的时候，panic 1将会传播给
			//外层被推迟执行的匿名函数，并与之关联。
			func() {
				panic(1) //与当前执行的匿名函数关联。
				//一旦panic 1 发生了，在这个新开启的goroutine中就同时存在了两个
				//panic。一个（panic 0）与入口函数关联，另一个(panic 1)则与当前执行
				//的匿名函数调用关联。
			}()
		}()
		panic(0) //与入口函数关联。
	}()
	select {} //阻塞等待
}

//用写入文件方式代替println测试多个goroutine的painc与recover，
//作为一种参考，知识点前面那已经介绍
func TestPanicRecoecrInMultiGoroutines2(t *testing.T) {
	const (
		mainFile string = "Main.txt"
		subFile1 string = "sub1.txt"
		subFile2 string = "sub2.txt"
	)
	goroutineProtecter := func(taskName, logFilePath string) {
		r := recover()
		if r != nil {
			message, ok := r.(string)
			if ok {
				writeFile("recover the panic: '"+message+"' in "+taskName+"\n", logFilePath)
			} else {
				writeFile("recover unkonw type panic in "+taskName+"\n", logFilePath)
			}
		}
	}
	writeFile("main task begin\n", mainFile)

	wg := sync.WaitGroup{}
	wg.Add(2) //为主任务增加两个并发执行的子任务。只有两个任务都完成，主任务才能继续工作。
	//子任务1
	go func() { // 子routine的入口函数。
		defer wg.Done()
		defer goroutineProtecter("sub taks 1 ", subFile1)

		func() {
			writeFile("sub task 1  executed\n", subFile1)
			panic("sub-task1-panic ")
		}()
	}()
	//子任务2，通过休眠3秒，模拟较重的写入任务
	go func() { // 子routine的入口函数。
		defer wg.Done()
		defer goroutineProtecter("sub taks 2", subFile2)
		func() {
			time.Sleep(3 * time.Second) //休眠3秒，模拟较重任务
			writeFile("sub task 2  executed\n", subFile2)
			panic("sub-task2-panic")
		}()
	}()
	wg.Wait() //等待并发执行的所有子任务完成。
	writeFile("main task end!\n", mainFile)
}

// 向指定路径的文件写入字符串内容
func writeFile(content, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(content)
	err = file.Sync()
	if err != nil {
		return err
	}
	return nil
}
