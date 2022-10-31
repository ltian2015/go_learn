/*
main包是包括可执行程序入口点函数文件的“特殊包”。
其特殊之处在于必须用main命名该包，同时该包中必有一个文件包括了可执行程序入口点函数,也就是main()函数。
main包被编译后，会形成一个可执行文件。
而普通包被编译后只会形成包文件(.a文件)
如果main包中有多个文件，那么必须用 go run "." 命令运行才能保证该包所有的文件被编译和运行。
本程序的main包中有两个文件，包括了入口main(）函数的main.go文件，以及模拟登录的login.go文件。
如果仅使用go run "main/package/path/main.go",则只会编译main.go文件，login.go不会被编译。
而main包中其他文件(login.go)无法被编译,从而导致“链接器”无法“链接”这些被main函数所在文件中使用，
但在其他未被编译文件(login.go)中所定义的变量，因此，当main包中有多个文件，
应当进入main包所在路径，使用go run "." 命令进行编译并运行，才能使main包中所有文件被编译。

在本程序的main包中，主要学习和展示包变量的定义与初始化顺序、依赖包的引入、
“空标识符”在依赖包的引入与变量声明中的用途，以及，初始化函数init()的特点与用途。

*
*/
package main

//com.example/golearn是在go.mod文件中定义的GO模块，“GO模块”代表了
//go.mod文件所在路径下的所有go package的集合。所以，go module相当于为这些构成这些
//go package文件的“相对根路径 /”定义了一个唯一标识名，知道这个唯一标识名就可以
//得到所有go package文件的“相对根路径”，import 将该“相对根路径/”与具体包文件相对
//"偏移路径"组合就可以访问到该包的文件。这样，其他依赖该包的文件就可以导入包。
//go语言的包名与所在目录名可以不一致，但是一个目录下的文件只能在一个包里，即：在同一个包中。
//go语言通过路径名引用包，但是注意，由于一个目录下的文件只能在一个包里，所以GO允许路径名与包路径不一致。
//比如，"com.example/golearn/reflectlearn"是路径名，但是该路径名go文件实际上在名为“myreflect”的包中。
//又比如，"com.example/golearn"是路径名，该路径下的包名是“main”。
//当然，由于一个路径下只有一个包，所以可以为导入路径所存储的包设置“别名”
// 形如： alias "import path",比如：myreflect "com.example/golearn/reflectlearn"
//最佳实践是路径名与包名一致。
//通常情况下，如果导入了包但不使用，就会无法通过编译。但如果用空标识符 _ 作为导入包的别名，则可以避免此问题。

import (
	//使用空标识符(blank identifier) _ 作为包的别名可以引入一个包而暂时不使用（预留未来使用）
	"fmt"
	"sync"
	"time"

	_ "com.example/golearn/usefunction"
)

// init()函数是一个特殊函数，用来完成包的初始化（变量初始化、验证与校验）。其特殊之处有两点：
// 1.自动执行。当包及所依赖包的变量都被初始化之后，就会自动调用init函数。
// 2. 该函数具有游动性质，可以在包中的一个文件或多个文件中多次按需出现。
// GO程序的初始化运行在一个单独的goroutine中(主goroutine中)，但是该goroutine可以创建其他的goroutines，
// 而这些goroutines可以并发运行。
// 如果package p 引入了（imports）package q, 则q的init函数的完成要发生在p的任何init函数之前。
// 所有的init函数的完成都“同步先于”main.main函数的启动。
// 这个初始化函数用来验证是否已经登录。登录操作在login.go文件的init()函数中完成。
func init() {

	if !IsLogin {

		panic("not login error")
	}
	var a string
	go func() {
		a = "hello"
		time.Sleep(1 * time.Second)
		println("hi,i'm in  init funcntion,but ,main function mabye over !")
	}()
	println(a)
}

var welcome string

// 本文件中第二个init()初始化函数，用来初始化welcome变量，即，使用login.go文件中定义的变量生成欢迎信息。
func init() {
	welcome = "hello " + UserName + " " + LoginTime.Format("2006-01-02 15:04:05")
}

// main包中的main()函数是特殊函数，它可执行程序的“入口点”。
// 注意，只有main包中的main()函数才是可执行程序的“入口点”。
// 其他非main包中也可以有main()函数，但非main包中的main()函数只是一个普通的函数。
// 作为入口点的main()函数所在的文件名未必一定命名为main.go,也可以是其他文件名，比如entrypoint.go
func main() {
	println("main function has started....")
	var noUse int
	// TODO: 使用空标识符解决已声明但是暂不使用（但未来打算使用）的变量无法通过编译的问题。
	// 主要用来写TODO的临时代码
	_ = noUse
	println(welcome)
	fmt.Printf("%b\n", 168)
	//functionalstyle.DataBandingTest()

	//VarCopy.TestVarCopy()
	//memorymanage.MemoryLayoutTest1()
	//contextlearn.Run()
	//learnSelect.DeadlockSinceChannelOpNeverOccured()
	//concurrentLearn.CopyForbidTest()
	//deferpanic.TestPanicDeferAndRecover()
	//EmptyStructEqualityTest()
	//EqualityAndCopy.TestCopyForbid()
	//controlflow.TestSwitch()
	//controlflow.TestIf()
	//cb.FixDeadLock()
	wg := sync.WaitGroup{}
	wg.Add(2) //增加两个任务。
	go func() {
		defer wg.Done()
		func() {
			//panic("panic in a goroutine")
		}()
	}()
	go func() {
		defer wg.Done()
		func() {
			//panic("panic in other goroutine")
		}()
	}()
	wg.Wait()
	time.Sleep(2 * time.Second)
	println("main function  exit")

	//	cb.ForRangeOnChannel()
	var lock sync.Mutex
	lock.Unlock()

}
