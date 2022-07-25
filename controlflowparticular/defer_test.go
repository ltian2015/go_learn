/**
defer、panic/recover是go语言中特殊的程序控制流。
这里主要关注GO语言的error机制、panic异常抛出机制、defer（延迟执行）机制，
以及defer与recover结合捕获并消除panic异常，实现程序从异常中恢复的机制。
本文件重点关注defer机制
**/

package controlflowparticular

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

/**
defer关键字是推迟函数执行程序流程执行机制。
defer机制与go routine可以看作go语言程序中一种特殊的控制流(control flow)。
defer语句的语法是： defer 函数调用表达式。
在计算机语言中，函数调用其实是一种表达式，与其他的表达式一样可以被求值，
”函数调用表达式“的求值就是执行函数，得到返回结果。
下面就是常见的defer 语句:
		defer  myFn(i)
		defer myFn(otherFunc(i))
		defer  obj.MemberFn（i）
		defer obj.MemberFn(i*2)
		defer obj.MemberFn(otherFn(i))
一个完整GO函数生命周期包括正常执行与退出两个阶段，只有当两个阶段都执行完毕，函数才会返回到上级函数。
函数的退出阶段会由return语句或抛出panic语句或runtime.Goexit函数的调用所引发。
defer语句的语义是：在外层函数正常阶段执行defer语句，而在外层函数的退出阶段执行”函数调用表达式“，
外层函数执行defer语句的主要操作是把 “函数调用表达式”压入defer-call堆栈之中，
在外层函数的退出阶段，如果defer-call堆栈为空，就什么都不做，外层函数会返回到上级函数。
否则就会按照先进后出的顺序，逐个弹出”函数调用表达式“，予以执行。

即使”函数调用表达式“执行时发生了panic，即使没有被”发生panic的函数“所recover，也不会影响外层的对
后续defer-call堆栈中的”函数表达式“的执行，因为，外层函数已经进入了退出阶段。 Go的panic机制下，
一个函数最多只能关联一个panic相关联，这种情况下，外部函数只能关联到最后发生的panic。

切记的一点就是，defer-call堆栈中存储的是“函数调用表达式”。

总之，当外层函数执行遇到defer “函数调用表示式” 语句时，go运行时并不会立即执行”函数调用表达式“，
而是把”函数调用表达式“所涉及的函数指针及调用该函数的参数都求值后压入defer-call堆栈中，
在函数进入退出阶段时，按照后进先出的顺序执行每个”函数调用表达式“。

因此，defer 函数调用有以下6个特点：
1. defer语句执行时，“函数调用表达式”中涉及的变量，包括被推迟的函数自身的函数指针、函数输入参数、函数接收者参数
会在defer语句执行时求值，然后形成完整的可调用的”函数调用表达式“压入defer-call堆栈中。
（1）”函数调用表达式“中所出现的变量值的后续变化对defer调用时的现场参数值没有影响。
（2）”函数调用表达式“中如果有变量计算表达式或另外另一个函数调用表达式作为被推迟调用的函数的参数，
    在defer语句被执行时，变量计算表达式会被立即求值，而另外的”函数调用表达式“也会被立执行，也就是说，会
	立即调用了另外的一个函数表达式，被把该函数的执行结果作为被推迟的函数调用表达式中的组成部分而被一同压栈，
	以便在退出阶段执行被推迟的函数调用。

2.若推迟执行的是闭包函数，若”闭包函数体“的代码中使用了外层函数或其他处定义的外部变量，
  在defer语句在执行时并不会对其求值。
  因此，被推迟的闭包函数中所捕获的外部变量在外层函数正常阶段执行defer语句时的值与在被推迟执行”函数调用表示“
  在退出阶段执行的值可能不一样，这个变量值可能会被外层函数或其他被推迟而先执行的函数所改变。
  这是容易引入bug的一点，需要谨慎对待。

3.defer函数调用能够改变主函数的返回结果。
因为defer 函数执行发生在主函数的退出阶段，这意味着主函数
仍在其生命周期之中，主函数的返回结果如果被命名，而这个命名的返回结果其实只是主函数生命周期中的一个局部
变量，能够在主函数的退出阶段被defer函数所改变。这种方式需要谨慎对待，否则容易出现Bug。

4.defer函数最终一定会被调用，无论主函数还是其他被defer的函数是否抛出了panic。
  调用的顺序为最先defer的函数最后调用。

5.有些系统函数不能被defer。
  被defer的函数在主函数退出阶段被自动调用，因此被defer的函数即使有返回结果也无法被外层函数的逻辑所用。
  因此，go会抛弃被推迟调用的函数的返回结果，因此，可以被推迟的函数的返回结果一定能够允许被抛弃，
  值得注意的是，系统内置函数(buidin包与unsafe包)中，除了copy与recover函数外，其他系统内置
  函数如果存在返回值，不允许被抛弃，因此，这些系统函数不能被defer。

6.如果被defer的函数是nil，尽管外层函数在正常执行阶段执行efer语句时不会抛出panic，
  但是，在外层函数在退出阶段执行该空函数的调用表达式时会抛出panic。

defer函数的用途主要有三点：
1. defer函数可以用于优雅地释放被主函数所申请的系统资源（比如，file，socket，lock等）。也就是在申请到资源后，立即defer 资源释放函数，
这样，函数正常执行完毕后，一定会释放资源。此场景下 ，defer相当于Java中的finally，但打开多个资源的时候，
go defer机制要比使用java finally机制更加优雅。

2.defer函数与recover函数配合，可以用于恢复（recover）主函数中收到的panic异常。
  如果一个函数抛出了异常，就会导致该函数进入自身的退出阶段，在自身的退出阶段通过recover函数
  来截获发生的异常，如果函数自身退出阶段没有处理panic异常，该panic异常就会成为上级的外层函数的panic，
  导致上级的外层函数进入退出阶段，在上级的外层函数退出阶段（的defer函数中）仍有机会使用recover函数截获
  panic异常，并进行处理，如果所有外层函数在退出阶段都没有处理，就会导致线程的入口函数抛出该panic，
  并使线程入口函数进入退出阶段，如果线程入口函数退出阶段没有捕获并recover该panic，就会导致
  该线程因存在panic而崩溃，go运行时就会由于任何一个线程的崩溃，而让整个程序崩溃。
  关于如何抛出panic与恢复panic的机制，详见panic_test.go。

3.利用defer 机制可以进行函数递归嵌套调用时对递归调用进行调用层次的追踪与分析。
**/

//TestDeferConctextCapture函数展示了推迟函数对推迟现场情况，反映了defer特点中的第1点和第2点。
//同时该函数展示了闭包函数对外层函数变量的捕获与defer对延迟调用时求值
func TestDeferConctextCapture(t *testing.T) {
	//part1 改变被推迟执行的函数变量本身
	var f = func() { //闭包函数1
		println(false)
	}
	defer f()    //f本质上是一个含有“函数指针”的复杂结构，此时被压入defer-stack栈的函数指针是闭包函数1，因此，最终打印false。
	f = func() { //闭包函数2。
		println(true)
	}
	var names = [5]string{"Kite", "Jone", "Mike", "Weiline", "Kisa"}
	for i := 0; i < 5; i++ {
		//part1 被推迟闭包函数捕获了调用参数现场，但是所捕获的外层函数变量值与推迟调用时的现场值不同步。
		defer func(index int) {
			name := names[index]
			println(i, "deferd say hello ", name) //i是被闭包函数捕获的变量，由于闭包函数被推到退出阶段执行，故而对捕获的变量i求值时，变量i的值已经变为了5.
		}(i) //这里的i会被defer时求值，i的值作为被推迟函数调用现场参数值，与函数指针一同被压入defer-call堆栈。

		//part2 同样的闭包函数，对比之下，未被推迟的闭包函数的所捕获外层函数变量同步求值。
		func(index int) {
			name := names[index]
			println(i, " hello ", name) //i是被闭包函数捕获的变量，由于闭包函数立即执行，所以对变量进行了立即求值，反映出与外部函数同步的i值。
		}(i)
	}
}

//TestModifyNamedResultByDefer函数展示了defer函数可以对外层函数命名的返回结果进行修改。
//这种方式要谨慎应用，但是否考虑为返回的抽象类型的结果插入装饰能力，或者无论出现什么
// 异常情况都让完成函数得到缺省的结果呢？有待于进一步研究。
func TestModifyNamedResultByDefer(t *testing.T) {
	mainFunc := func(welcomeWord string, name string) (welcomeStatement string) {
		defer func() {
			welcomeStatement += ", thanks" //此处，闭包函数捕获了外层函数的返回结果变量，并在外层函数退出阶段执行了对返回结果值的修改。
		}()
		return welcomeWord + " " + name
	}
	println(mainFunc("hello", "susan"))

}

//该TestDeferNilFunc函数演示了defer空函数的效果。
//该案例说明在defer中使用函数变量是一个坏习惯。
func TestDeferNilFunc(t *testing.T) {
	println("begin")
	var f func() //函数类型变量f的零值是nil
	//下面的defer语句在压栈前对函数变量f进行了求值，求值结果为nil，因此将nil函数压入到defer-call stack中。
	//这会导致在函数退出阶段执行对nil函数的调用，因此会抛出panic。
	defer f() //nil 函数被压入defer-call 栈
	//此时为f变量赋予一个非空函数值，已经于事无补
	f = func() {
		println("hello")
	}
	//此语句会被正常执行
	println("i reach the end")
	return //此后，会抛出nil函数调用的异常。
}

//TestElgentResuorceRelease函数演示了使用defer方式，在一段使用两个资源的业务逻辑中优雅地释放两个资源。
func TestElgentResuorceRelease(t *testing.T) {
	const fileName = "log.txt"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic("can not open file ./log.text")
	}
	defer file.Close() //释放文件资源
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock() //释放同步锁
	println("will  write text to log.txt")
	file.WriteString("it's ok,if you see this line " + time.Now().Format("2006-01-02 15:04:05.999999999Z07") + "\n")
	//file.WriteString("it's ok,if you see this line " + time.Now().Format("“2006-01-02 15:04:05.999999999Z07:00”") + "\n")
}

//deferCloseMultiFilesInBadWay展示了一个函数可能一次性打开大量资源，
//如果使用defer机制在在该函数退出时释放资源，会导致资源不合理占用，与性能的下降。
func deferCloseMultiFilesInBadWay(paths []string) error {
	const fileContent string = "hello"
	for _, path := range paths {
		file, err := os.Open(path) //打开文件资源
		if err != nil {
			return err
		}
		defer file.Close() //推迟到主函数结束阶段关闭文件
		_, err = file.WriteString(fileContent)
		if err != nil {
			return err
		}
		err = file.Sync() //将文件的内存提交到永久存储中。
		if err != nil {
			return err
		}
	}
	return nil
}

//deferCloseMultiFilesInGoodWay同样展示了需要打开大量文件，并采用defer机制优雅关闭文件的
//一种比较好的方式，主要是用匿名函数包括对每一个文件的打开与关闭处理，这样会及时关闭打开的资源。
func deferCloseMultiFilesInGoodWay(paths []string) error {
	const fileContent string = "hello"
	for _, path := range paths {
		//在匿名函数中使用defer机制完成单个文件的打开....处理....关闭操作
		//这样，每次匿名函数调用结束时，被打开的文件资源就会被释放。
		if err := func() error {
			file, err := os.Open(path) //打开文件资源
			if err != nil {
				return err
			}
			defer file.Close() //推迟到主函数结束阶段关闭文件
			_, err = file.WriteString(fileContent)
			if err != nil {
				return err
			}
			err = file.Sync() //将文件的内存缓存内容提交到永久存储中。
			if err != nil {
				return err
			}
			return nil
		}(); err != nil {
			return err
		}
	}
	return nil
}

//------------------defer 特殊用例1：通过refer机制配合高阶函数重启出现panic的goroutine-----------
func TestRestartPanicGoroutine(t *testing.T) {

	go stubbornTaskExcutor("task1", businessWorkFunc)
	go stubbornTaskExcutor("task2", businessWorkFunc)
	select {} //阻塞主gotoutine，从而可以显示两个子goroutine的执行情况

}

//stubbornTaskExcutor函数定义一个了执着的任务执行者，
//这是一个高阶函数，它会执行传入的工作函数，如果传入工作函数在执行过程中出现了panic，
//它会重新启动一个goroutine继续运行自身，直至任务完成，所以称之为执着的stubborn任务执行者。
//通过本用例的实现机制，我们可以设计一个更加完善的，可以重试一定（配置）次数后再最终报错的类似于akka任务执行框架。
func stubbornTaskExcutor(taskName string, workFunc func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(taskName, " is carshed. restart it now !")
			//重启线程意味着每panic一次，就会在当前goroutine下重新启动一个新goroutine。
			//其实，也可以不用重启那么多协程，即：去掉go关键字。
			go stubbornTaskExcutor(taskName, workFunc)
		}
	}()
	workFunc()
}

//	businessWorkFunc函数模拟一个具体的业务处理函数。
func businessWorkFunc() {
	//模拟一个工作负载
	fmt.Println("workd begin")
	time.Sleep(time.Second)
	//模拟一个随机的非业务处理逻辑范围内的意外问题的发生
	if time.Now().UnixNano()&0x3 == 0 {
		panic("unexcepeted situation occur!")
	}
	fmt.Println("workd end")
}

//

//------------------defer 特殊用例2：通过refer与panic机制的配合实现跨多层函数嵌套调用的转移跳转-----------
//必须强调的是，这是一种不易读的风格，只有极特殊的场景可以使用。慎用！
func TestRemoteJump(t *testing.T) {

	n := func() (result int) {
		//恢复错误
		defer func() {
			//通过defer/recover 机制获得n层函数调用直接抛出的结果。
			if r := recover(); r != nil {
				if num, ok := r.(int); ok {
					result = num
				}
			}
		}()

		//多层嵌套调用，会有一层调用中产生panic
		func() {
			//嵌套调用
			func() {
				//嵌套调用
				func() {
					//嵌套调佣
					func() {
						panic(123) //通过panic，结束多层嵌套调用的函数，直接将结果返回给捕获panic的函数。
					}()
				}()

			}()

		}()
		return
	}()

	fmt.Println("n=", n)
}

//------------------defer 特殊用例3：通过refer与panic实现更加简洁的错误检查-----------
//对于一个顺序执行多个操作操作步骤的程序，每一个操作步骤都存在三种情况：
//1.出现错误，不能继续。
//2.没有错误，但是逻辑决定不能继续。
//3.继续执行。
//如果采用常规方式定义执行步骤，则执行步骤的定义为：
//  func stepN()(continue bool,err error)
//执行步骤的主程序必须判断两个返回结果，这样主程序就比较啰嗦.
//而采用panic和recover则可以简化处理。
func TestSimplifyErrorCheck(t *testing.T) {
	var simulateStep = func(stepNum int, haveError, continueNextStep bool) {

		if haveError { //发生错误，后续步骤不能继续执行
			panic("error ourrue in step " + strconv.Itoa(stepNum))
		} else {
			if !continueNextStep {
				panic(nil) // 没有错误，但是业务逻辑决定后续步骤不能继续执行。
			}
			fmt.Println("task over in step ", stepNum)
			//没有错误，后续步骤可以继续执行
		}

	}
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()
	//简洁地实现了错误检查。
	simulateStep(1, false, true) //模拟没有错误，继续执行
	simulateStep(2, false, true) //模拟没有错误，继续执行
	simulateStep(3, false, true) //模拟没有错误，继续执行
	simulateStep(4, false, true) //模拟没有错误，继续执行
	simulateStep(5, true, false) // 模拟发生错误，肯定不能继续执行的步骤，第二个参数被忽略。
	simulateStep(6, false, true) //得不到执行了
}

//------------------defer 特殊用例4：利用Defer机制实现嵌套函数的调用追踪-----------
func TestTracing(t *testing.T) {
	tracer := NewTracer("")
	childFn := func() {
		//注意，tracer.Trace("childFn")会在childFn函数正常执行阶段被执行求值，其执行结果作为
		//tracer.Untrace调用表达式的参数在childFn函数的退出阶段执行。
		//在childFn函数退出所执行的函数调用表达式形如：tracer.Untrace(”tracer.Trace执行结果")
		defer tracer.Untrace(tracer.Trace("childFn"))
	}
	parentFn := func() {
		//注意，tracer.Trace("childFn")会在parentFn函数正常执行阶段被执行求值，其执行结果作为
		//tracer.Untrace调用表达式的参数在parentFn函数的退出阶段执行。
		//在childFn函数退出所执行的函数调用表达式形如：tracer.Untrace(”tracer.Trace执行结果")
		defer tracer.Untrace(tracer.Trace("parentFn"))
		childFn()
	}
	//注意，tracer.Trace("childFn")会在TestTracing函数正常执行阶段被执行求值，其执行结果作为
	//tracer.Untrace调用表达式的参数在TestTracing函数的退出阶段执行。
	//在childFn函数退出所执行的函数调用表达式形如：tracer.Untrace(”tracer.Trace执行结果")
	defer tracer.Untrace(tracer.Trace("TestTracing"))
	parentFn()
}
func TestPrinter(t *testing.T) {
	printer := func(message string) string {
		println(message)
		return "END " + message
	}
	defer printer(printer("TestPrinter"))
	println("\thandling in TestPrinter")
}
