package functions

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"
)

// !!! 这个例子展示了闭包的基本性质
// !!! 在Go语言中，闭包（closure）是一种函数值(可以翻译为字面函数或函数字面量)，
// !!! 它可以引用其函数体之外的变量(是字面函数出现的程序上下文，一般是父函数体)。
// !!! 闭包函数可以访问并修改这些外部变量，这些变量在闭包(字面函数)中被绑定，
// !!! 它们的生命周期会延续到闭包不再被使用为止。
func TestClosure1(t *testing.T) {
	//incrementor是一个函数b变量， 等价于func incrementor() func() int {......}
	incrementor := func() func() int {
		i := 0              //!!! 注意，此变量是将被下面闭包函数所捕捉的变量，其生命周期将延续到闭包（字面函数）不在使用（调用）为止。
		return func() int { //这是一个字面函数，也叫函数值，也叫闭包。
			i++ //!!!被闭包捕获的外部变量，只要闭包被调用，其生命周期就不会结束。
			return i
		}
	}

	incrementtorA := incrementor() //得到一个闭包（函数值）
	println(incrementtorA())       //使用闭包，也就是调用闭包值，闭包捕获变量的值增加1，打印 1
	println(incrementtorA())       //使用闭包，，闭包捕获变量的值再增加1，打印 2
	println(incrementtorA())       //使用闭包，，闭包捕获变量的值再增加1，打印 3
	incrementtorB := incrementor() //重新获得另外一个闭包（函数值）,该闭包所捕获的外部变量与上一个闭包所捕获的外部变量不是同一个实例,且初始值被该闭包所初始化。
	println(incrementtorB())       //使用闭包 ，闭包捕获变量的值增加1，打印 1
	println(incrementtorB())       //使用闭包 ，闭包捕获变量的值增加1，打印 2
	println(incrementtorA())       //使用闭包，，闭包捕获变量的值再增加1，打印 4
}

// ------------------------------Student packae-----------------
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

//------------------------Teacher package
// !!! 可以将该函数看作是“别的包（package）”中的函数，即，它访问不到任何有关Student的信息。
// !!! 也就是它不依赖Student所在的包。但是老师可以在阅卷之后为学生加分，只要将学生的加分操作参数传递给他即可。

func TeacherIncrementScore(studentScoreIncrementor func(deata int) int, delat int) int {
	return studentScoreIncrementor(delat)
}

// !!! 用途1：展示了闭包用于封闭和持久化内部状态，并提供对封闭状态的操作
// !!! 使得在外部可以操作对象内部状态，而无需知道对象是谁。
// !!! 这个测试函数模拟了引用了两独立包：Student 包和Teacher包，并且让Teacher为学生加分。
// !!! 这里就利用了闭包（函数值）的状态封装能力将学生与老师解耦。
// !!! 主要思路就是分数状态持有方要（1）将状态被闭包所捕获（2）对外提供闭包函数。

func TestClosure2(t *testing.T) {
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

// ------------------------------ package 函数能力增强  function Capability enhance  -----------------
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

func TestClosure3(t *testing.T) {
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

// ----------------------package 延迟计算(加载)--------
func lazyCmopute(a int, b int) func() int {
	result := 0
	return func() int {
		if result == 0 { // !!! 只有第一次调用时才计算结果
			time.Sleep(1 * time.Second) // 模拟计算延迟
			result = a * b
			return result
		} else {
			return result
		}
	}
}

func TestClosure4(t *testing.T) {
	// !!! 这个例子展示了闭包用于延迟计算
	// !!! 通过闭包来封装一个计算函数，只有第一次调用时才进行计算
	compute := lazyCmopute(3, 5)
	println(compute()) // 输出: 15 (经过1秒延迟)
	println(compute()) // 输出: 15 (没有延迟)
}

// ----------------------package 作为工厂方法按已有模板创建可变的对象实例----------------
// !!! 闭包中捕获了用作模板的其他变量，然后生成只需要少量参数的工厂函数。
// !!! filename=name + suffix ，创建用给定后缀（suffix）作为模板固定参数的文件名称生成器。
func fileNameGenerator(suffix string) func(name string) string {
	// 这里可以写校验suffix的逻辑
	//
	return func(name string) string {
		return name + suffix
	}
}

func TestClosure5(t *testing.T) {
	// !!! 这个例子展示了闭包作为工厂方法
	// !!! 通过闭包来生成具有特定后缀的文件名
	txtFileNameGenerator := fileNameGenerator(".txt")
	docFileNameGenerator := fileNameGenerator(".doc")
	println(txtFileNameGenerator("file1")) // 输出: file1.txt
	println(txtFileNameGenerator("file2")) // 输出: file2.txt
	println(docFileNameGenerator("file3")) // 输出: file1.doc
	println(docFileNameGenerator("file4")) // 输出: file2.doc
}

// ---------------------------------package 作为共享变量在并发goroutine中的拷贝----------------
// !!! 这个例子展示了闭包用于在并发goroutine中共享变量
func TestClosure6(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(10) // !!! 10个goroutine
	for i := 1; i < 11; i++ {
		go func(i int) {
			defer wg.Done()
			// !!! 这里的i是闭包捕获的变量，每个goroutine都有自己的i值的拷贝（如果被捕获的变量是指针就不秒了！）
			index := i
			i = i * 10 // !!! i值被拷贝，所以可以被goroutine自由修改，不要担心
			time.Sleep(time.Duration(i) * time.Millisecond)
			println("Goroutine", index, "completed after", i, "milliseconds")
		}(i) // !!! 注意，这里传递了i的值，避免了闭包捕获同一个变量的问题

	}
	wg.Wait() // 等待所有goroutine完成
}

//
