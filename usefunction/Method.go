package usefunction

/**
函数（行为）与数据绑定，此时，函数就是一种带有“接收者（reciever）”的特殊函数，
这种函数也称之为接收器的方法，即：method。这就是GO语言中对面向对象编程的支持方式。
函数与数据的绑定方式有两种：
1. 绑定到数据引用，也就是绑定指针，就是把指向数据的指针作为函数的接收器（receiver）。
2. 绑定到数据值，就是以数据值作为函数的接受者(receiver).
GO语言中，行为与数据的两种绑定方式可以分别用来实现对象的成员方法和纯函数的不可变编程。
需要注意的是：
   1.函数只能与同一个包中的数据类型进行绑定, 不能绑定其他包中的数据类型。
		如果允许PackageA 中的函数（行为）绑定了 PackageB中的数据类型，那么，如果想要
		通过PackageB中数据类型的数据实例调用其行为函数（在PackageA中），就需要同时引入
		PackageA和PackageB。而只有同一个包中的行为与数据类型进行绑定，就不会带来额外的包引入。
		所以绑定到数据类型上的行为只应是针对该类型自己的内聚式闭合行为。

也就是这些行为函数所需要的输入与输出参数类型，都是被绑定数据类型自身，或构成与被绑定数据类型
相关的类型，或Go语言系统公用的基础类型，总之，尽量不要在行为中引入其他非GO系统包的数据类型。

这两种方式主要是与指针的概念与用途相关，在下文中有详细探讨。
https://medium.com/@annapeterson89/whats-the-point-of-golang-pointers-everything-you-need-to-know-ac5e40581d4d
这里，对文章中的重要内容进行了翻译，如下：
如果有人问你为什么要在GO中用指针，你就可以告诉他们，“指针效率高，因为在GO中你所传递的每个东西
都是值（value），所以，指针让我们传递的是数据所在的地址而不是数据的值，通过将普通变量与指针变量分开，
可以避免无意间的数据变更，当我们想要在其他函数中改变一个值的时候，使用指针的方式就可以让我们访问到实际值，
而不是值的拷贝”

上面的道理告诉了我们如何去使用指针？主要有以下用途：
1. 如果我们想要一个函数去改变接收者（receiver）的某些数据，我们必须使用指针来作为函数的接受者（receiver），此时，函数就是这个对象的成员方法。
2. 如果我们不想要一个函数去改变接收者的某些数据，那么就用数据的值而不是数据的指针来作为函数的接受者（receiver）。此时，函数就是一个纯函数。
3. 如果我们有一个非常大的数据，我们通过传递对这个数据所在位置的引用值可以获得更高的性能，也就是指针。
4. 由于通过指针可以改变数据，在并发（Concurrency）操作的时候，就会导致共享数据的不同步，
   因此，指针容易带来并发安全问题。所以在GO routines之间要避免使用指针来共享数据，尽量不要在channel中传递
   指针值或引用值。

使用方法需要注意的是：
1.!!!*****如果符合某个接口形式的函数绑定所接收者是某个结构体的指针类型，
  !!! 那么，在go中这个结构体的指针类型是一种新的类型，编译器认为是这个结构体的指针类型实现了接口，
  !!! 而不是结构体类型实现了该接口。

2. 无论方法的接收者类型被定义为普通变量还是指针变量，在调用该方法时，
   使用普通变量还是指针变量都可以调用该方法，go编译器自动完成了普通变量与指针变量之间的转换。
**/

type Student struct {
	ID   int
	Name string
}

// 结构体的成员方法（Method Member），因为拷贝传递过来的是指针，通过指针操作会改变所指向的结构体的内容。
// 注意：即使用的不是指针而是变量来调用该函数，也不会改变它是成员方法的语义，
// GO编译器会自动完成接收器变量与指针之间的转换。
// 此时，GO会先自动取变量的地址创建指针，然后再用指针作为参数调用该方法,比如：
// var s Student=Student{ID:1,name:"lantian"}
// s.changeNameMehod(”ltian“)  //等价于
// (&s).s.changeNameMehod(”ltian“)
// GO 自动取地址，生成指针进行调用。
// s,函数的接收者(receiver)，这里是个指针，意味着，本函数是一个改变接受者数据的成员方法。
func (s *Student) changeNameMethod(name string) {
	s.Name = name
}

// 纯函数风格，拷贝传递过来的Student结构体，并命名为s，改变并返回拷贝后的结构体，
// 该函数不会对原结构体有任何的改变，所以是纯函数风格。
// 注意：即使使用的是指针而不是结构体变量来调用该函数，也不会改变它是纯函数的语义。
// 此时GO会先用指针找到所指向的Student变量，拷贝Student变量，然后再用该拷贝的Student对象调用该方法，比如：
// var s Student=Student{ID:1,name:"lantian"}
// sp:=&s
// sp.changeNameMehod(”ltian“)  //等价于
//
//	(*sp).changeNameMehod(”ltian“)
//
// s,函数的接收者(receiver)，这里是一个数据值，意味着，本函数是一个不改变接收这数据的纯函数。
func (s Student) changeNameFunction(name string) Student {
	s.Name = name
	return s
}

// 测试结构体的成员方法与纯函数风格。
func DataBandingTest() {
	var s1 Student = Student{ID: 1, Name: "lantian"}
	var s2 Student = Student{ID: 2, Name: "lan lan"}
	println("s1 : Student{", s1.ID, ",", s1.Name, "} @ ", &s1)
	println("s2 : Student{", s2.ID, ",", s2.Name, "} @ ", &s2)
	sp1 := &s1
	sp2 := &s2
	sp1.changeNameMethod("ltian1") //等价于下面
	s1.changeNameMethod("ltian2")
	s2_1 := sp2.changeNameFunction("lanlan2") //等价于
	s2_2 := s2.changeNameFunction("lanlan2")  //等价于
	println("s1 : Student{", s1.ID, ",", s1.Name, "} @ ", &s1)
	println("s2 : Student{", s2.ID, ",", s2.Name, "} @ ", &s2)
	println("s2_1 : Student{", s2_1.ID, ",", s2_1.Name, "} @ ", &s2_1)
	println("s2_2 : Student{", s2_2.ID, ",", s2_2.Name, "} @ ", &s2_2)
}

// 函数可以绑定到处理基础类型和未命名类型之外的任何数据类型上，包括整数类型、字符串类型以及另一个函数类型。
type MyInt int

// 函数绑定到指针类型的数据上，注意，不是绑定到MyInt上，MyInt和*MyInt不是同一个类型。如果
// isPositive() 函数是某个接口的唯一方法，那么是 *MyInt类型实现该接口，而不是MyInt实现了该接口。
func (i *MyInt) isPositive() bool {
	return int(*i) > 0
}

type StringHandler func(s string) string

// 绑定到了函数类型
func (sh *StringHandler) printHandleResult(s string) {
	result := (*sh)(s)
	println(result)
}
