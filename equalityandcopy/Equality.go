/**
equalityandcopy包主要包括学习如何判断两个变量是否相等的代码，以及变量是否允许拷贝（通过赋值）的检查与控制。
**/

package equalityandcopy

/**
这个文件主要用来展示如何正确进行变量相等性的判断。
在任何语言中，判断两个变量所代表的数据是否相等都是一件很重要的事情，因为很多时候，
相等性判断用来决定逻辑是否执行以及如何执行。相等性是引起bug的常见来源之一。
比如，在很多if语句都需要相等性判断。而Map的key也要求必须是”可比较的“数据类型。

关于相等性判断，在GO语言存在如下法则,详见GO语言规范关于各种变量的比较规则：https://golang.org/ref/spec#Comparison_operators
1. 首先==或!=操作符只能用于类型相同的变量，即使二者底层类型相同也不能比较，会产生编译错误。
2. 要看被比较的两个变量的数据类型是否是“可比较的（comparable）”数据类型，只有可比较的数据类型才能直接用
==或!=操作符来判断是否相等或不等。否则对于不可比较的类型使用==或!=操作符会产生编译错误。
如下数据类型是“可比较”的数据类型
boolean
numeric    //注意，NaN（Not a Number,超出表达范围的数）与NaN不相等 math.NaN()==math.NaN()求值结果为false
string
pointer
channel
interface types
structs 如果该struct所有字段的类型都是“可比较的”类型。
array 如果该数组的元素类型是“可比较的”类型

而如下3种数据类型是“不可比较”的数据类型，使用“==”或“!=” 来判断是否相等或不等，会导致编译错误。
Slice
Map
Function
3.对Struct和Array来说，如果Struct没有Field（或0个Field），Array没有元素（或0个元素），
这样的类型的变量被称作“0内存（zero-size）变量”，“O内存变量”具有相同的内存地址。但不具比较性。
4.对于Slice、Map和不可比较的struct，可以用reflect.DeepEqual进行比较，
  DeepEqual比较的逻辑是递归方式比较所有内部成员（包括公开和私有成员）的相等性。
5.如果觉得对GO本身提供的变量比较逻辑不满意，希望采用自定义的方式进行比较，可使用
google/go-cmp 包进行自定义的比较。该包导入路径为： "github.com/google/go-cmp/cmp"
更多关于比较的详细例子见：
https://medium.com/golangspec/equality-in-golang-ff44da79b7f1
https://golangbyexample.com/struct-equality-golang/
https://www.geeksforgeeks.org/how-to-compare-equality-of-struct-slice-and-map-in-golang/
**/

import (
	"fmt"
	"unsafe"
)

// 测试底层类型相同的不同类型的变量相等性
func UnderlyingTypeEqualityTest() {
	type MyInt int
	var mi MyInt = 5
	var i int = 5
	println(mi, i)
	//GO不同类型的变量不能用==操作符进行比较，即使二者底层类型相同也不行，必须要进行类型转换，才能比较。
	//以下语句无法通过编译。
	//	var isEqual bool = (mi == i)
}

// 下面函数测试可比较的struct
func ComparableStructEqualityTest() {
	type ComparableStruct struct {
		id      int               //可比较类型
		name    string            //可比较类型
		isMale  bool              //可比较类型
		weight  float32           //可比较类型
		channel chan string       //channel是可比较类型
		parent  *ComparableStruct //指针是可比较类型
	}
	chanel1 := make(chan string) //channel是引用类型的变量，引用的背后变量不是一个（引用类型的地址指针field值不同），二者就不相等。
	//ch1:=make(chan string),ch2:=make(chan string), ch1和 ch2分别引用了不同的背后变量，所以二者不相等。

	var1 := ComparableStruct{id: 1, name: "lantian", isMale: true, weight: 65.5, parent: nil,
		channel: chanel1}
	var2 := ComparableStruct{id: 1, name: "lantian", isMale: true, weight: 65.5, parent: nil,
		channel: chanel1}
	if var1 == var2 {
		println("var1 and  var2 are equal") //这个将会是打印结果
	} else {
		println("var3 and  var4 are not equal")
	}

	//在两个结构中，虽然结构体的channel域各自都是通过make(chan string)调用产生的，表面看起来似乎一样，但是两个结构体的channel域引用了两个不同的背后变量，是不同的channel。
	var3 := ComparableStruct{id: 1, name: "lantian", isMale: true, weight: 65.5, parent: nil,
		channel: make(chan string)}
	var4 := ComparableStruct{id: 1, name: "lantian", isMale: true, weight: 65.5, parent: nil,
		channel: make(chan string)}
	if var3 == var4 {
		println("var1 and  var2 are equal")
	} else {
		println("var1 and  var2 are not equal") //这个将会是打印结果
	}

	//下面定义了一个不可比较的结构体类型
	type NotComparableStruct struct {
		id      int                   //可比较类型
		name    string                //可比较类型
		chidren []NotComparableStruct //不可比较类型-Slice类型。
	}
	ncs1 := NotComparableStruct{1, "lantian", nil}
	ncs2 := NotComparableStruct{1, "lantian", nil}
	println(ncs1.id, ncs2.id)
	//以下的代码无法通过编译，因为NotComparableStruct是不可比较的结构体类型，该类型变量不能作为==操作符的操作数。
	//if ncs1 == ncs2 {
	//	println("ncs1 and ncs2 are equl")
	//}
}

// 测试指针变量的相等性
// 无论指针指向的类型是可比较还是不可比较的，所有类型的指针都是可比较的。
// 指针类型相同，且指向同一个内存地址，则二者是相同的指针。值为nil的同类指针都相等
// 指向两个不同的 zero-size 变量的指针可能相同，也可能不同。这是GO语言规范中的说，原因大概
func PointerEqualityTest() {
	var p1, p2 *string
	name := "foo"
	p1 = &name
	p2 = &name
	println("p1==p2 ? ", p1 == p2)    //true，值同为 nil
	println("*p1==*p2 ?", *p1 == *p2) //true
	slice1 := []int{1, 2, 3}          //Slice是不可比较的变量类型
	var ps1, ps2 *[]int
	ps1 = &slice1
	ps2 = &slice1
	println("ps1==ps2 ? ", ps1 == ps2) //true，无论何种类型的指针，指针类型本身是比较的类型，只要指向同一个地址，就是相同的指针。
	//println("*ps1==*ps2 ? ", *ps1 == *ps2)  编译无法通过，因为指针指向的变量类型是不可比较的。
}

// 比较通道类型对象的相等性,通道对象是GO基本的并发原语（ concurrency primitive），判断两个通道对象是否相等是否重要。
// 如果两个通道对象都是nil,或者两个通道对象都是由“同一次make方法调用”所产生的。
// 原理：通道底层实现是用一个包含了阻塞队列的指针字段（Field）来引用阻塞队列的结构体，也就是所为的“引用类型”。
// 所以，该结构体的所有字段（Field）必须相同，那么channel才会相同。
func ChannelEqualityTest() {
	var ch1, ch2 chan int
	println("ch1==ch2 ? ", ch1 == ch2) //true,因为ch1和ch2都是nil
	ch1 = make(chan int)
	ch2 = make(chan int)
	println("ch1==ch2 ? ", ch1 == ch2) //false,因为二者背后引用的底层阻塞队列不同（表面上看是由两次make调用产生的）。
	ch2 = ch1                          //引用拷贝，使得表达两个引用变量底层的两个结构体各字段（Field）值相同。
	println("ch1 memory address is ", &ch1)
	println("ch2 memory address is ", &ch2)
	println("ch1==ch2 ? ", ch1 == ch2) //true，二者相同，因为背后引用的底层阻塞队列相同。
	comp := func(cha chan int, chb *chan int) {
		println("cha==chb ?", cha == *chb)
	}
	comp(ch1, &ch2) //true，二者相同，因为背后引用的底层阻塞队列相同
}

// 空对象，或0内存对象的比较。
// 空对象，或0内存对象是指没有元素的数组，或没有字段的结构体。
// 注意，所有空对象或0内存对象的地址都相同，共享同一个内存地址，
//但是按照GO语言规范中有关指针变量比较的说明“指向不同的空对象指针作比较，得出来的结果可能相同，也可能不同”。
//大概是因为Go与编译器智能决定变量是在堆栈中还是堆中。由于在堆中的对象要进行垃圾回收，垃圾回收的CPU优先级高于
//goroutine的CPU优先级，所以会降低GO程序执行效率。为了减少堆的垃圾回收，GO编译器发现如果变量没有被两个以上
//函数所共用，那么就在栈中为其开辟内存，否则在堆中开辟内存。如果在堆中开辟内存，就存在动态回收的情况，
//比较动态变化的指针地址没有意义，因此会返回false。有关GO 内存分配规则见：
//https://medium.com/faun/golang-escape-analysis-reduce-pressure-on-gc-6bde1891d625

func EmptyObjEqualityTest() {
	type EmptyStruct struct{}
	type EmptyIntArray [0]int
	type EmptyStringArray [0]string
	esa := EmptyStruct{}
	esb := EmptyStruct{}
	println("empty struct var esa address is ", &esa)
	println("empty struct var esb address is ", &esb)
	isValueEqual := esa == esb
	println("esa == esb ? ", isValueEqual) //true，空结构体的值都相同。所有类型相同的空对像都像相等。因为都是一个实例。
	var eai EmptyIntArray                  //空的整数数组
	var eas EmptyStringArray               //空的字符串数组
	println("empty int array address is ", &eai)
	println("empty string array address is ", &eas)
	//println("eai == eas ? ", eai == eas) //类型不同，不允许比较
	//println("&eai == &eas ? ", &eai == &eas) //指针类型不同，不允许比较

	//将不同类型的空对象的指针（指针地址都相同）转换成统一的指针unsafe.Pointer进行比较，结果为fals。
	//这个结果体现了规范中所说的空对象（0内存对象）指针比较没有意义，结果不确定。
	// unsafe.Pointer 本质上是一个指向int的指针类型,其定义如下：
	// type ArbitraryType int
	// type Pointer *ArbitraryType
	println("unsafe.Pointer(&esa) ==uasfe.Poinnter(&eai) ? ", unsafe.Pointer(&esa) == unsafe.Pointer(&eai)) //false，空对象指针比较
	println("unsafe.Pointer(&esa) ==uasfe.Pointer(&esb) ? ", unsafe.Pointer(&esa) == unsafe.Pointer(&esb))  //false，空对象指针比较

	//比较不同类型的空对象的内存地址
	/**
	uintptr 是一个可以存储指针地址的整数类型，而不是指针类型，所以不可以通过*uintptr的方式来访问指针地址所指向的对象。
	uintptr可以进行算术运算，而unsafe.Pointer则不能进行整数运算，因此，uintptr主要配合unsafe.Pointer来使用，
	通过把unsafe.Pointer转换成uintptr，进行地址的算术运算后，再转换成unsafe.Pointer进行对象的操作。
	unsafe.Pointer可以表示任何类型的指针，任何类型的指针可以转换为unsafe.Pointer，同时，unsafe.Pointer可以转换
	为任何类型的指针。
	*/
	println("&esa ==&esb ? ", uintptr(unsafe.Pointer(&esa)) == uintptr(unsafe.Pointer(&eai))) //true，所有空对象内存地址相同

}

/**值拷贝发生在赋值的时候。显式的赋值就是赋值语句，隐式的变量赋值包括：函数调用时的参数传递，
   函数的返回语句将结果赋值给结果变量。
   无论何种类型的变量，在GO语言中，所有的赋值都是值的拷贝。
   注意：在Go语言中，如果一个值的类型是T，而有一些方法与*T绑定（关联），那么对T值拷贝就比较危险。
// In general, do not copy a value of type T if its methods are associated with the pointer type, *T.
**/
// 下面的函数就是用来验证GO语言的赋值是值拷贝。
func TestVarCopy() {
	type StructObj struct {
		Id int
	}
	var originStructObj StructObj = StructObj{1}
	var copyStructObj = originStructObj
	copyStructObj.Id = 22
	fmt.Printf("origin struct id is : %v \ncopy struct id is : %v \n", originStructObj.Id, copyStructObj.Id)
	if originStructObj.Id != copyStructObj.Id {
		fmt.Println("GO 语言中，结构（struct）类型的变量的赋值是值拷贝，创建了一个新的结构体值，而不是引用拷")
	}
	/**
		   在GO语言中，有三个所谓的“引用类型“，即 ：slice, map，chan。
		   本质上，每个引用类型就是用包含了底层数据内存指针的struct作为数据类型来操作底层数据。
		   比如,slice就是这样的一个数据类型定义方式（runtime/slice.go）：
		   type slice struct {
		        array unsafe.Pointer  //sliece底层数组的内存地址指针
		        len   int
		        cap   int
	        }
			所以，slice是对数组的引用，map是对哈希表的引用，chan是对阻塞队列的引用。
			其他两个引用类型的实现定义见：runtime/map.go 和 runtime/chan.go。
			由于引用类型本质上是struct类型，所以，引用类型的变量的赋值其实是创建了一个新的struct实例（该实例有自己的内存和不同的内存地址），
			该新实例的值拷贝了赋值来源struct变量的所有数据，当然也包括了底层的内存地址。
			因此，两个引用变量实际上还是操作同一块内存。
	**/
	//originSlice实际上是一个struct变量，包含了底层数组的指针field。
	var originSlice = []int{1, 2, 3, 4, 5}
	//copySlice实际上是另一个struct变量，底层数组的指针field的值拷贝自originSlice，
	//所以二者实际上操作的是同一个内存地址。
	var copySlice = originSlice
	copySlice[0] = 100
	fmt.Printf("origin slice[0] is : %v \ncopy slice[0] is : %v \n", originSlice[0], copySlice[0])
	fmt.Printf("origin slice address is : %p \ncopy slice address is : %p \n", &originSlice, &copySlice)

}

// GO语言开发者认为，在go中，表示多个数据域复合对象struct拷贝会产生潜在的问题（主要是浅拷贝所带来的问题）。
func TestStructCopyProblem() {
	//表示所有“教师实体”的数据类型
	type teacher struct {
		id   int
		name string
	}
	//表示所有“学生实体”的数据类型
	type student struct {
		id      int
		name    string
		teacher *teacher //指向了一个老师实例的指针。
	}
	t1 := teacher{id: 1, name: "teacher liu"}
	s1 := student{id: 1, name: "student li", teacher: &t1}

	s2 := s1 //s2拷贝了s1,由于 teacher是一个指针变量，所以，会导致指针变量的值的拷贝，使两个学生导师老师都是一个。
	// 第一个潜在问题，浅拷贝可能导致数据共享，即：指针数据类型变量的拷贝导致指向同一个数据。
	//在某种场合下，这是希望的结果，但在有些场合，就不合适。
	// 比如,一个学生与导师是一对一的关系时，这种拷贝就导致了学生的导师并未拷贝出一个新教师对象。
	s2.teacher.name = "teacher mao" //通过学生2的指针更改导师名字，其实与学生1共享了一个导师。
	fmt.Println(s1.teacher.name)    //会打印出teacher mao，这是由于浅拷贝导致的数据共享。
	//
	studentMap := make(map[student]string)
	studentMap[s1] = s1.name
	if name, ok := studentMap[s2]; ok {
		println("find student of s2 ,name is ", name)
	}
	if name, ok := studentMap[s1]; ok {
		println("find student of s1,name is ", name)
	}
}
