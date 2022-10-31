/*
*
escapeanalysis包给出了GO内存的分配机制以及内存逃逸的原因分析，以供高性能程序开发者参考。
在编程语言的编译优化原理中，分析指针动态范围的方法称之为逃逸分析。
需要让GO在运行时分配内存的语句（命令）只有变量声明、new或make等语句，而变量内存要么分配在
goroutine\thread堆栈内存中，要么分配在多个goroutine\thread共享的堆内存中（Java用new创建的变量都会在堆内存中。）
在支持“指针”概念的编程语言中，通过指针引变量时可能会发生变量存活期超过所在函数（或闭包，或LOOP块)
存活期而使用的可能。但这种情况发生时，由于变量所在函数调用结束时，该函数在goroutine堆栈中对应的堆栈桢（Stack Frame）就会消失，
如果在goroutine堆栈中为变量分配内存，在函数存活期之外再通过指针访问该变量，就会导致内存访问错误。
为了避免这种情况发生，Go编译器首先会在编译期间检查是否存在变量存活期超过所在函数存活期而使用的情况，
如果存在这种情况，就把内存分配实现代码编译为在堆中为变量分配内存。然后运行时再通过垃圾回收器进行垃圾回收，从而保证不会有内存泄漏。
但是垃圾回收例程的优先级很高，在堆内存为大量的变量分配内存，会导致垃圾回收的工作量很大，
因此使得正常业务代码无法高效利用CPU而降低程序性能。对于时间关键性任务，因垃圾回收导致大幅性能下降的详细介绍见：
https://medium.com/eureka-engineering/understanding-allocations-in-go-stack-heap-memory-9a2631b5035d

对于处理时间要求很高的服务程序，要尽量避免变量逃逸到堆中的情况，以减少垃圾回收的时间。
事实上，GO语言是在编译时决定在堆栈中还是堆中为变量分配内存。主要通过堆根据源代码生成的抽象语法树进行变量
数据来源的链路回溯，结合链路中变量取址与指针取值操作，来决定变量的内存分配位置。详见：
因此，GO语言编译器在编译源代码时就能够清楚地知道变量应该在何处分配内存。
所以，GO编译器提供了一个工具可以检查程序源代码来发现哪些变量逃逸到了堆内存中。
在go build命令打开 -gcflags开关就可以堆源代码进行逃逸分析，逃逸分析工具以源代码包为单位进行分析，例子如下：
go build -gcflags '-m' "com.example/golearn/memorymanage/escapeanalysis"

	或输出更详细说明的命令
	 go build -gcflags '-m -m' "com.example/golearn/memorymanage/escapeanalysis"

上述命令对escapeanalysis包中的源代码文件进行逃逸分析。
通常，使GO编译器在堆内存中为变量分配内存，也就是变量逃逸到堆中的情况主要有以下几种：
1包级变量不一定会引起内存逃逸（包级变量在goroutine的堆栈中，函数在变量在堆栈的函数桢中），

	除非在包级声明了指针变量指向了包级以下程序中范围内的变量。

2.函数输出（返回了）其内部声明的局部变量的指针给外部——调用该函数的更大函数。
3.在for循环体外部声明的指针变量指向了for循环体内部声明的局部变量。

	for循环都会导致在goroutine堆栈中创建一个“堆栈桢”，就像函数调用时创建一个“函数堆栈桢”一样。

4.闭包外部声明的指针变量指向了闭包内部声明的局部变量。
5.变量占用内存过大，超出了堆栈的内存范围（1<<20）。

总之，只要指针所指的变量在小于指针所在的范围内声明或创建（new），就会导致变量逃逸到堆中。
指针所指变量在大于或等于指针所在范围内声明或创建（new）的变量时，就不会导致变量逃逸到内存中。

*
*/
package escapeanalysis

var PackageVar int = 20                //不会逃逸到堆中，会分配在主存的
var PackagePointer *int                //在本程序中引起了逃逸，指向了包级以下的（函数级）范围内的变量
var PackagePointer2 *int = &PackageVar //在本程序中没有引起逃逸。指针指向同范围变量，就不会导致该变量逃逸到堆中。
type DoNoEscape interface {
}




type MyInt32 int32

func (mi *MyInt32) Do() {
	referencedByPackagePointer()
}

// 验证包级指针导致函数中变量逃逸到堆中的情况。
func referencedByPackagePointer() {
	y := PackageVar
	y++
	PackagePointer = &y // 会导致y逃逸到堆中。因为PackagePointer是包级指针，而y是函数级变量。
}

// 验证函数返回指针变量导致局部变量逃逸到堆中的情况。
func returnPointer() *int {
	var myint int = 100
	var name string = "hello"
	var strPointer *string = &name //不会导致name逃逸到堆中，指针指向的是同级范围变量。也未超出范围拷贝指针。
	*strPointer += " veryone"
	return &myint //会导致myint 逃逸到堆中。函数返回内部局部变量的地址指针变量给外部调用者。
}

// 验证闭包内部变量逃逸到堆中的情况。
func closureVarEscape() {
	var intPointer *int //本程序中，该指针引起了闭包中变量的逃逸。
	//下面是一个函数字面量，也就是匿名函数，或闭包。
	func(i int) {
		var doubleI int = i * 2
		var sum int = doubleI + i
		intPointer = &sum //会导致sum逃逸到堆中。闭包范围之外的指针指向了闭包内部的局部变量。
	}(5)
	*intPointer++
}

// 验证for循环体内变量逃逸情况
func loopEscape() {
	var intPointer *int
	for i := 0; i < 1; i++ {
		var doubleI int = i * 2
		intPointer = &doubleI //会导致doubleI逃逸到堆中，因为定义在循环体内部的变量被循环体外部指针引用。
	}
	*intPointer++
}

// 验证变量过大导致goroutine堆栈存放不下而逃逸到堆中的情况。
func bigVar() {
	nums := [2000000]int{} //因为数组变量内存太大，超过了goroutine堆栈的最大值减去其他变量已预定的堆栈内存，所以逃逸到了堆中。
	nums[0] = 1
}

// 验证if块外指针指向if块内变量是否会导致变量逃逸到堆中的情况，结果表明，不会导致内存逃逸。
func ifScopeEscape() {
	var i int = 2
	var pointInt *int //不会导致if语句块内的变量逃逸。
	if i > 1 {
		var varInIfScop int = 10
		pointInt = &varInIfScop
		*pointInt++
	}

}
