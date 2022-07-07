package equalityandcopy

/**
这个文件主要展示限制对象拷贝的措施
 GO语言开发者认为，在go中，表示多个数据域复合对象struct拷贝会产生潜在的问题（主要是浅拷贝所带来的问题）。
 所以，在有些GO 对象禁止拷贝（复制，函数调用传参，返回值接收），也就是，一旦源
 对象被使用后，使用源对象的拷贝对象就会发生Panic异常。
 同时利用go vet 工具，可以在编译阶段检查出不允许拷贝的对象，比如：
 go vet "com.example/golearn/equalityandcopy"
 或更为详细的信息输出：
 go vet  -v "com.example/golearn/equalityandcopy"
 注意，有些IDE工具会自动编译并进行vet检查，并显示编译错误和vet检查发出警告信息，比如vscode。
 不允许拷贝的实现分为个步骤：
 1.让对象成为不应被拷贝的对象，不应被拷贝的对象，一旦在源代码中发生了赋值拷贝，
   则在编译阶段可以被 go vet工具检查出来。但，不应拷贝对象即使发生了赋值拷贝，也可以通过编译，
   而且在运行期间不会引发Panic异常样。
 2.在不应被拷贝的对象的方法中添加“拷贝检查程序”，在运行期间，检查被拷贝后，抛出panic异常。

 不应被拷贝的对象只要实现了   sync.Locker 接口，一但对此类对象进行赋值拷贝，
 就可以被  go  vet 工具所检查出来。
 sync.Locker 接口声明如下：
type Locker interface {
	Lock()
	Unlock()
}
任何struct 只要实现上述两个空方法，一旦被赋值拷贝，即可被go vet工具检测出来。

**/

import (
	"sync/atomic"
	"unsafe"
)

//
//GO的Go语言中如何实现一个不应被拷贝struct呢？
//下面这个struct 用两个空方法实现了  sync.Locker接口。所以一旦发生了赋值拷贝，就可以被go vet工具检查出来
//
type MyNoCopyStruct struct {
	Name string
}

func (*MyNoCopyStruct) Lock()   {}
func (*MyNoCopyStruct) Unlock() {}

/**
    只要把NoCopy放在任何一个结构体内部，这个对象就可以快速成为一个不应被拷贝的对象，
	一旦被赋值拷贝，都会被go vet 工具能够检查出来。
*/
type ShouldNoCopy struct{}

func (*ShouldNoCopy) Lock()   {}
func (*ShouldNoCopy) Unlock() {}

/**
使用 NoCopy 快速实现了一个不应被拷贝对象。
*/
type QuickShouldNOCopy struct {
	nocopy ShouldNoCopy
	Name   string
}

// copyChecker 存储了copyChecker对象自身的地址，用来进行拷贝检查。
// 当一个copyChecker对象被创建之后，第一次使用之前，所存储的自身地址都是初始化的零值，
//当该对象的check()方法被第一次调用（使用）后，才会被赋值为自身的地址。
// 在第一次被使用之后（自身地址变量指向了自己的地址），再拷贝这个对象，调用check（）方法就会
//检查出其所存储的自身地址变量值与调用check()方法对象的地址不一致，从而检查出对象发生了拷贝。
type CopyChecker uintptr

func (c *CopyChecker) check() {
	if uintptr(*c) != uintptr(unsafe.Pointer(c)) && // 首次检查存储自身地址的变量值与对象指针的地址是否相同 。
		!atomic.CompareAndSwapUintptr((*uintptr)(c), 0, uintptr(unsafe.Pointer(c))) && //如果存储自身地址的变量值是初始的零值，则更换为实际的地址。使用了CompareAndSwapUintptr可以保证操作的原子性，如果交换成功，交换函数返回true，否则交换函数返回false。
		uintptr(*c) != uintptr(unsafe.Pointer(c)) { //再次检查所存储自身地址的变量值与对象指针的地址是否相同。
		panic("sync.Cond is copied")
	}
}

//使用ShouldNoCopy和CopyChecker快速实现一个对象实例不可拷贝的类型。
//而且该类型的每个成员方法在代码之前都应调用checker的check()方法，
//检查对象是否发生过拷贝。所以，不可拷贝对象不应暴露除了方法之外的成员数据。
//这个例子暴露了不可拷贝的对象的成员数据，严格来讲，是不合适的。
//在系统库中，sync.Cond类型的对象实例就是不可拷贝的对象。
type QuickNOCopy struct {
	nocopy  ShouldNoCopy //不应拷贝的对象
	checker CopyChecker  //拷贝检查
	Name    string
}

func (qnc *QuickNOCopy) PrintName() {
	qnc.checker.check() //检查是否发生过拷贝
	println(qnc.Name)
}
func (qnc *QuickNOCopy) SayHello() {
	qnc.checker.check() //检查是否法发生过拷贝
	println("hello everyone!")
}

//测试不应别拷贝和不可被拷贝对象的代码。
func TestCopyForbid() {
	var shouldNocopy MyNoCopyStruct = MyNoCopyStruct{"lari"}
	copy1 := shouldNocopy //该拷贝赋值语句，编译期间可被vet工具检查出来，但不会引起运行异常。
	println(copy1.Name)
	qShoulNocopy := QuickShouldNOCopy{Name: "lari"}
	otherQsnocp := qShoulNocopy //该拷贝赋值语句，编译期间可被vet 工具检查出来，但不会引起运行异常。
	println(otherQsnocp.Name)
	var quickNOCopy QuickNOCopy = QuickNOCopy{Name: "lantian"}
	quickNOCopy.PrintName() //第一次使用
	theCopy := quickNOCopy  //该拷贝赋值语句，编译期间可被vet 工具检查出来，但不会引起运行异常。
	theCopy.SayHello()      //拷贝了不可拷贝的对象，会引起运行时异常
}
