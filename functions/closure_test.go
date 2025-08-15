package functions

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

// !!! 下面两个例子展示了闭包的基本性质
// !!! 在Go语言中，闭包（closure）是一种函数值(可以翻译为字面函数或函数字面量)，
// !!! 在代码编写时，闭包的代码中可以使用字面函数代码所出现的程序上下文中定义的变量（通过变量的标识符），
// !!! 在代码运行时，闭包会捕获所在上下文代码运行时所产生的的外部变量实例。
// !!! 尽管这些被捕获的变量实例所在的上下文代码（往往是闭包代码出现的上级函数代码）可能已经执行完毕，
// !!! 但这被捕获的变量实例的生命周期会延续到闭包不再被使用为止。
// !!! 如果把闭包所在的上下文函数称为“环境函数”，被闭包所捕获的“环境函数”中的状态变量，
// !!! 这些状态变量被"环境函数所初始化"，然后被一个或者多个闭包函数捕捉并多次操作，就相当于一种“隐式的面向对象”。
// !!! 操作“隐式的面向对象”内部状态的各闭包函数甚至不在同一个goroutine中，它们可以通过通道或者所进行协作。
func TestCaptureVariable1ByClosure1(t *testing.T) {
	//incrementor是一个函数b变量， 等价于func incrementor() func() int {......}
	incrementorFactory := func() func() int {
		i := 0              //!!! 注意，此变量是将被下面闭包函数所捕捉的变量，其生命周期将延续到闭包（字面函数）不在使用（调用）为止。
		return func() int { //这是一个字面函数，也叫函数值，也叫闭包。
			i++ //!!!被闭包捕获的外部变量，只要闭包被调用，其生命周期就不会结束。
			return i
		}
	}
	incrementtorA := incrementorFactory() //!!!incrementorFactory() 函数在本次执行过程中才会有变量i的实例存在，闭包就是捕获就是incrementorFactory在该次执行中所产生的变量i的实例。
	println(incrementtorA())              //使用闭包，也就是调用闭包值，闭包捕获变量的值增加1，打印 1
	println(incrementtorA())              //使用闭包，，闭包捕获变量的值再增加1，打印 2
	println(incrementtorA())              //使用闭包，，闭包捕获变量的值再增加1，打印 3
	incrementtorB := incrementorFactory() //!!!incrementorFactory() 函数又一次运行，闭包捕获的incrementorFactory变量i的实例与上次运行所捕获的不是同一个实例。
	println(incrementtorB())              //使用闭包 ，闭包捕获变量的值增加1，打印 1
	println(incrementtorB())              //使用闭包 ，闭包捕获变量的值增加1，打印 2
	println(incrementtorA())              //使用闭包，，闭包捕获变量的值再增加1，打印 4
}

// !!! 这个程序表明两点：
// !!! 1. 产生闭包函数的主函数运行完毕并退出后，被闭包所捕获的主函数变量仍然还在。
// !!! 2. 主函数一次运行所生成的变量实例，可以被多个闭包捕获，它们捕获的都是同一个变量实例。
func TestCaptureVariable1ByClosure2(t *testing.T) {
	CreateDoublerAndGetter := func(i int) (doubler func(), getter func() int) {
		defer func() {
			println("CreateDoublerAndGetter函数已经退出")
		}()
		doubler = func() {
			i = i * 2
		}
		getter = func() int {
			return i
		}
		return doubler, getter //!!!double和get两个标识符不过是内存中代码段的地址，是闭包函数的入口地址
	}
	double, get := CreateDoublerAndGetter(10) //!!!double和get两个标识符不过是内存中代码段的地址，是闭包函数的入口地址
	println(get())                            //打印10
	double()
	println(get()) //打印20
	double()
	println(get()) //打印40
}

// !!! 这个例子展示了多层嵌套的闭包可以捕获最外层函数的变量，这里是msg
func TestClosure6(t *testing.T) {
	wg := &sync.WaitGroup{}
	msg := "hello,"
	for i := 1; i < 11; i++ {
		wg.Go(func() { //第一层闭包
			func(a int) { //第二层闭包
				// !!! 这里的i是闭包捕获的变量，每个goroutine都有自己的i值的拷贝（如果被捕获的变量是指针就不秒了！）
				index := a
				a = a * 10 // !!! i值被拷贝，所以可以被goroutine自由修改，不要担心
				//time.Sleep(time.Duration(i) * time.Millisecond)
				println(msg, " Goroutine", index, "completed after", a, "milliseconds")
			}(i)
		}) // !!! 注意，这里传递了i的值，避免了闭包捕获同一个变量的问题
	}
	wg.Wait() // 等待所有goroutine完成
}

func TestClosureWithGoroutine(t *testing.T) {
	var wg sync.WaitGroup
	msg := "hello"
	// !!! 这两个在闭包中的线程对所捕获的msg变量的改变顺序是不确定，无法预料。
	wg.Go(func() {
		defer fmt.Println("f1 over!")
		msg = "welcome"
	})
	wg.Go(func() {
		defer fmt.Println("f2 over!")
		msg = "hi"
	})
	wg.Wait()
	fmt.Println(msg) //!!! 可能是hi，也可能是 welcome，取决于f1,f2哪个最后执行
}

func TestClosureWithGoroutine2(t *testing.T) {
	var wg sync.WaitGroup
	age := 1
	wg.Add(2)
	go func() { //这个闭包捕获的是主函数一次运行轮次中age变量.
		defer wg.Done()
		age = 100
	}()
	go func() { //这个闭与上个闭包捕获age变量与上个闭包所捕获的age变量在主函数同一个运行轮次中，都是同一个age变量.
		defer wg.Done()
		age = 200
	}()
	wg.Wait()
	fmt.Println(age)
}

// /////////////////////////////////////////////////////////////////////
// !!! 用途1：闭包用于封闭和持久化内部状态，并提供对封闭状态的操作
// !!! 使得在外部可以操作对象内部状态，而无需知道对象是谁。
// !!! 这个测试函数模拟了引用了两独立包：Student 包和Teacher包，并且让Teacher为学生加分。
// !!! 这里就利用了闭包（函数值）的状态封装能力将学生与老师解耦,使得老师看不到学生是谁，但仍然可以为学生的卷子加分。
// !!! 主要思路就是分数状态持有方要（1）将状态被闭包所捕获（2）对外提供闭包函数。
// ------------------------------模拟 tudent packae-----------------
type student struct {
	id    int
	name  string
	score int
}

// 模拟对学生分数的增加，
func (st student) ScoreIncrementSimulater() func(deata int) int {

	return func(delta int) int {
		st.score = st.score + delta //!!! 此闭包内捕获的变量不只是st.core，还有整个st
		return st.score
	}
}
func (st *student) ScoreIncrementor() func(deata int) int {

	return func(delta int) int {
		st.score = st.score + delta //!!! 此闭包内捕获的变量不只是st.core，还有整个st
		return st.score
	}
}
func (st *student) ScoreIncrementor2(delta int) int {
	st.score = st.score + delta
	return st.score
}

//------------------------ 模拟 Teacher package
// !!! 可以将该函数看作是“别的包（package）”中的函数，即，它访问不到任何有关Student的信息。
// !!! 也就是它不依赖Student所在的包。但是老师可以在阅卷之后为学生加分，只要将学生的加分操作参数传递给他即可。

func TeacherIncrementScore(studentScoreIncrementor func(deata int) int, delat int) int {
	return studentScoreIncrementor(delat)
}

func TestClosureUseCase1(t *testing.T) {
	stduent := &student{id: 1, name: "lantian", score: 0}
	ScoreIncrementSimulater := stduent.ScoreIncrementSimulater()
	println(ScoreIncrementSimulater(15)) // 输出: hello lantian
	println(ScoreIncrementSimulater(20))
	println(stduent.score) // 输出: welcome lantian
	// !!!对外提供可以操作状态（这里闭包封装的状态是student的状态，被操作修改的score属性）闭包函数
	scoreIncrementor := stduent.ScoreIncrementor()
	//!!! 教师通过学生提供的状态操作方法来操作学生的分数
	TeacherIncrementScore(scoreIncrementor, 75)
	TeacherIncrementScore(scoreIncrementor, 25)
	// !!! 这种方式也可以，闭包的好处是减少了对student实例的传递。
	// !!! 因为，student实例的状态已被封装者在其闭包之内,可以减少不可控的修改和命名冲突
	TeacherIncrementScore(stduent.ScoreIncrementor2, 25)
	println(stduent.score) // 输出: welcome lantian

}

/////////////////////////////////////////////////////////////////////////////////////////////////
////////////// !!!闭包用途2： 函数能力增强  function Capability enhance  -----------------
//!!! 这个例子展示了闭包（字面函数）用于创建对给定函数的修饰函数，以增强给定函数所没有的能力

type RequestHandler func(request string, ctx context.Context) string

// !!! 输入给定的函数，返回增强后（增加过滤能力）的函数，注意，二者的形式一样
func RequestFilter(handler RequestHandler) RequestHandler {
	return func(request string, ctx context.Context) string {
		// !!! 这里可以添加一些请求过滤逻辑
		if request == "" {
			return "Invalid request"
		}
		if request == "ping" {
			return "pong"
		}
		return handler(request, ctx)
	}
}

// !!! 输入给定的函数，返回增强后（增加日志能力）的函数，注意，二者的形式一样
func RequestLogger(handler RequestHandler) RequestHandler {
	return func(request string, ctx context.Context) string {
		// !!! 这里可以添加一些请求日志逻辑
		log.Println(request, time.Now().Format(time.RFC3339))
		return handler(request, ctx)
	}
}

func TestClosureUsecase2(t *testing.T) {
	// !!! 这个例子展示了闭包的高级用法
	// !!! 通过闭包来封装请求处理器的过滤和日志功能
	handler := func(request string, ctx context.Context) string {
		return "Request processed: " + request
	}
	// !!! 使用闭包来添加过滤和日志功能
	filteredHandler := RequestFilter(handler)
	loggedHandler := RequestLogger(filteredHandler)

	ctx := context.Background()
	println(loggedHandler("ping", ctx))  // 输出: pong
	println(loggedHandler("hello", ctx)) // 输出: Request processed: hello
}

// ////////////////////////////////////////////////////////////////////////////////
// ///////!!! 闭包用途3：通过返回闭包函数实现延迟计算(加载)--------///////////////////////
// 实现方式：函数本身并不进行计算，而是定义了结果变量，通过返回的闭包函数进行计算，
// 该闭包函数捕获了宿主函数的中的结果变量，如果结果变量没有计算，那就进行计算，否则就直接返回捕获的结果。
// 实现原理：和面向对象中的，用对象内部状态来暂存计算结果一样。
// 闭包捕获的变量相当于在分配在“堆”上的，只有宿主函数与闭包函数共同可见的长期存在的内部状态。
//

func lazyCmopute(a int, b int) func() int {
	result := 0 //!!!被闭包所捕获的变量
	ok := false //!!!被闭包所捕获的变量
	return func() int {
		if !ok { // !!! 只有第一次调用时才计算结果
			time.Sleep(3 * time.Second) // 模拟计算负载
			result = a * b
			ok = true
			return result
		} else {
			return result
		}
	}
}

func TestClosureUseCase3(t *testing.T) {
	// !!! 这个例子展示了闭包用于延迟计算
	// !!! 通过闭包来封装一个计算函数，只有第一次调用时才进行计算
	compute := lazyCmopute(3, 5)
	println(compute()) // 输出: 15 (经过1秒延迟)
	println(compute()) // 输出: 15 (没有延迟)
	compute2 := lazyCmopute(0, 5)
	println(compute2()) // 输出: 15 (经过1秒延迟)
	println(compute2()) // 输出: 15 (没有延迟)
}

// ////////////////////////////////////////////////////////////////////////////////
// /////!!! 闭包通途4：实现模板工厂///////////////////////////////////////////////
// !!! 返回的闭包函数以给定参数作为模板，创建可变的对象实例
// !!! 闭包中捕获了用作模板的对象构成元素，只需要少量参数就能生产目标对象
// !!! filename=name + suffix ，创建用给定后缀（suffix）作为模板固定参数的文件名称生成器。
func fileNameGenerator(suffix string) func(name string) string {
	// 这里可以写校验suffix的逻辑
	//
	return func(name string) string {
		return name + suffix
	}
}

func TestClosureUseCase4(t *testing.T) {
	// !!! 这个例子展示了闭包作为工厂方法
	// !!! 通过闭包来生成具有特定后缀的文件名
	txtFileNameGenerator := fileNameGenerator(".txt")
	docFileNameGenerator := fileNameGenerator(".doc")
	println(txtFileNameGenerator("file1")) // 输出: file1.txt
	println(txtFileNameGenerator("file2")) // 输出: file2.txt
	println(docFileNameGenerator("file3")) // 输出: file1.doc
	println(docFileNameGenerator("file4")) // 输出: file2.doc
}

// ////////////////////////////////////////////////////////////////////////////
// !!!闭包用途5：通过暴露的闭包函数控（不可见的）制线程的运行。
// !!! 通过对外暴露的闭包函数控制隐藏的，在独立线程中执行的闭包函数的运行。
// !!! 这是一个典型的“隐式对象（Iterator）”，next、cancel 和push三个闭包函数
// !!! 共享相同的状态（入参s、Iterator_ch、Iterator_done、Iterator_canceld ）
// !!! 这个“隐式对象（Iterator）”的公开方法是next和cancel，而push则是私有方法，并在
// !!! “隐式对象（Iterator）”初始化时（GetIterator函数运行时）就在一个独立的线程中开始执行了。

func GetIterator[T any](s []T) (next func() (T, bool), cancel func()) {
	Iterator_ch := make(chan T)          //!!! 被多个闭包捕获的变量是“隐式对象”中共享的状态，命名尽量突出一下
	Iterator_done := make(chan struct{}) //!!! 被多个闭包捕获的变量是“隐式对象”中共享的状态，命名尽量突出一下
	Iterator_canceld := false            //关闭一个已关闭的通道会引发panic，所以用次标志避免取消函数多次调用导致的通道关闭异常
	cancel = func() {
		if Iterator_canceld {
			return
		}
		Iterator_canceld = true
		//!!! 用做取消通知的通道，使用关闭比发送信号更稳妥，
		//!!! 因为对取消信号处理线程可能因为正常结束而不能再读取取消信号通道，导致发送取消信号的线程阻塞
		close(Iterator_done)
		// Iterator_done<-struct{}{} //!!! 如果取消信号处理线程正常结束，那么运行发送取消信号的线程就会阻塞在这里
	}
	next = func() (T, bool) {
		result, ok := <-Iterator_ch
		return result, ok
	}
	push := func() { //!!! 该线程的执行进度被next、cancel两个闭包函数通过对共享通道的读取与关闭操作所控制
		for _, t := range s {
			select {
			case <-Iterator_done:
				close(Iterator_ch) //关闭数据通道
				println("关闭数据传送了")
				return //结束线程
			case Iterator_ch <- t:
			}
		}
		println("全部遍历，即将关闭数据传送通道")
		close(Iterator_ch) //如果全部读取，则关闭数据通道
	}
	go push()
	return next, cancel
}
func TestClosureUseCase5(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	next, cancel := GetIterator(s)
	defer cancel()
	for {
		data, ok := next()
		if !ok {
			break
		} else {
			println(data)
		}
	}
	data, ok := next()
	println(data, ok)
	println("---------------------------")
	next, cancel = GetIterator(s)
	println(next())
	println(next())
	cancel() //测试多次取消（但实际仅执行一次channel关闭操作），没有panic产生
	cancel()
	println(next()) //测试取消之后，再也读不出有效数据
	println(next())
}
