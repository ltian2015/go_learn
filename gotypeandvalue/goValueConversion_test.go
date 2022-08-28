package gotypeandvalue

import (
	"fmt"
	"testing"
)

/*****************************************************************************
值的类型转换的本质是希望能够运用另一个类型的方法集中的方法计算（处理）当前的值，也就是对“相同或兼容内存格式与布局”
切换不同的操作方法集，从而得到希望的结果。因此，一个类型的值如果可以转换另一个类型，那么前提就是目标类型的内存布局
一定能够兼容变量的内存格式与布局。但是为了防止内存布局兼容的类型之间出现不经意的类型转换，从而导致方法集的
滥用，GO不允许内存布局兼容的两类型之间进行隐式（直接）的类型转换，必须使用显式（强制）的类型转换，从而提高了类型安全。

1. 值的类型转换语法。
  1.1 编译器保证的安全类型转换——显式（强制）的静态类型转化。 如果一个值v，可以转换为某个类型T，那么转换的语法就是(T)(v). 尤其是，如果T“已定义的类型(defined type)”，
   那么往往简写为T（v）,称为“简化类型转换写法”。而(T)(v)可以称为“正规类型转换写法”。
   由于存在“简化类型转换写法”，所以， “正规类型转换写法”往往用于将值转换为“未经定义的类型”，
   因为未经定义类型是多个类型组成的复合类型，也就是字面类型，因此被()括起来的阅读效果更好。
   如果一个值x可以被隐式地转换为类型T  ，也就是说值x不需要类型转换语法就能够直接赋值给类型T的变量，
   那么，一定也可以显式地将其转换未类型T的变量。即：如果编译器允许 var t T=x，那么一定允许 var t =T(x)
  1.2 编译器无法完全保证安全的类型转换——运行时显式的动态类型转换——为不确定具体类型的接口值的断定“精准的”具体类型来源。
         ConcreteTypeValue,boolValue=InterfaceValue.(ConcreteType)
	  这个语法的语义是在运行期，判断不确定具体类型的接口值InterfaceValue其来源的真实类型是否为指定的
	  具体类型ConcreteType，如果是，就将该接口值转换为ConcreteType，并返回 (ConcreteTypeValue,true)，
	  否则返回（零值，false）。
	  之所以称其为无法完全保证类型安全，是因为必须依靠程序员检查转换结果。
******************************************************************************/
//--------------------------value convertion 语法练习 开始-----------------------//
type DefinedType struct {
	id   string
	name string
}
type Teacher struct {
	id   string
	name string
}

var data1 = DefinedType{id: "001", name: "xiaofei"}

// 将已定义类型 DefinedType的值data1转换为 未定义类型
//
//	struct {
//		id   string
//		name string
//	}
//
// 的值，采用“正规类型转换写法(T)(v)”的风格比较好，如下：
var data2 = (struct {
	id   string
	name string
})(data1)

// 写成简写风格有些不易读,如下:
var data3 = struct {
	id   string
	name string
}(data1)

// Teacher是已定义类型，采用简化写法T(v)进行类型转换比较好
var data4 = Teacher(data1) //T(v)风格
// 不安全的类型转换
func TestUnsafeTypeConvert(t *testing.T) {
	type Student struct {
		id   string
		name string
	}
	var i int = 99
	student, ok := (interface{})(i).(Student)
	if ok {
		println(student.id)
	} else {
		println("can not convert")
	}

}

//--------------------------value convertion 语法练习 结束------------------------//
/************************************************************************************
2.值的类型隐式（直接）转换规则 ：
	鉴于值的类型转换的本质是为内存格式或布局相同的数据切换方法集（Method  Set）。
	但是为了避免内存布局兼容的变量之间的发生无意地类型转换（隐式转换）而造成的误操作，GO制定了
	严格地类型转换规则，只有当同一类型具有不同类型名称是，才允许隐式的转换。
	这是考虑到go语言中可以为通过 type alias=T 来为类型定义别名，这就存在名称不一样但是实际上
	是同一个类型的情况。因此，如果两个类型本质上是同一种类型，那么可以隐式转换（直接赋值，var t T=x）。

***/
func TestSameTypeConvert(t *testing.T) {
	var aByte byte = 'a' //byte 是uint8的别名，二者是同一类型
	var i8 uint8 = aByte //直接赋值(内存拷贝)，隐式转换！
	var runeData []rune = []rune{1, 2, 3, 5}
	f := func(data []int32) {
		fmt.Println(data)
	}
	f(runeData) // rune是int32的别名，所以[]rune等同于[]int32,所以可以直接赋值（内存拷贝），隐式转换！
	fmt.Println(i8)
	//自定义类型与其别名
	type People struct {
		pId   string
		pName string
	}
	type PeopleTypeAlias = People                       //PeopleTypeAlias与People本质上是同一个类型。
	var p = PeopleTypeAlias{pId: "1001", pName: "Kite"} //直接赋值，隐式转换！
	fmt.Println(p.pId)

}

/***
2.2 与底层类型相关的转换规则。对与Tx类型的值x和"非接口类型"T，
    （1）如果二者（Tx和T）共享相同的底层类型（忽略struct中的元数据tag），因为，具有相同的内部布局，那么必须
	    将值x以显式的方式转换为类型T的值，这样才可以避免类型的误操作。
		因为，通过显式的类型转换，让程序员明确地知道切换了不同的（操作）方法集。
	（2）如果二者（Tx 和T）之间有一个是未经定义的类型（字面类型），并且二者的底层类型相同（要考虑struct中的元数据
	    tag，，那么x就可以隐式地转换为类型T。规则解析：因为未经定义的类型（字面类型）没有方法集（方法集为空，空既包含一切——作为
		源类型，可以转换成有限方法集，也什么都没有——作为目标类型，无方法可用，比较安全），两个类型值之间不存在方法集的差异。
		所以，类型转换一定是类型操作安全的，但是考虑到struct的tag元数据定义在类型上，通过类型的反射操作对元数据处理不同，
		因此，对于tag元数据不同，但其他定义都相同的struct类型，直接隐式转换会造成无意的错误，因为，不允许隐式转换，
		但是可以显式转换。
	（3）如果二者（Tx和T）的底层类型不同，但是 ，如果Tx和T都是未定义的指针类型（指针字面类型），并且，指向的
	    数据类型有相同的底层类型（忽略struct中的元数据tag），那么，必须将指针值x以显式的方式转换为类型T，才能使用T的方法集
		进行操作。 规则解析：对于两个未定义的指针类型（指针字面类型）Tx和T，指针类型所指向的底层数据内存布局一样，
		这样为相同内存布局换另一套方法集（指针字面类型的方法集取决于于所指向类型的方法定义），就不会造成内存错误，只不过是为了避免不经意的逻辑错误，需要程序员
		作显式（明显的）类型转换。但是以未定义的指针类型为底层类型的已定义类型（defined type）类型，由于
		还可以定义自己的方法集，为了避免误操作（方法的乱用），必须显示地从其底层的指针类型转换过来。
***/
//------------------------与底层类型相关的转换规则的练习，begin----------------------------------//
func TestUnderlyinngTypeSameConvert(t *testing.T) {
	type IntSlice []int
	type MySlice []int
	//以下三个类型的底层类型相同，都是组合类型，或称字面类型 或称未定义类型 []int
	var undefinedTypeS []int //未定义类型变量
	var is IntSlice          //已定义类型变量
	var ms MySlice           //已定义类型变量
	//is = ms  //编译错误，不能隐式转换,防止方法集不经意的误用。
	//ms=is // 编译错误，不能隐式转换,防止方法集不经意的误用。
	is = undefinedTypeS //可以直接的隐式转换，因为不存在方法集误用的可能。
	undefinedTypeS = is //可以直接的隐式转换，因为不存在方法集误用的可能。
	ms = undefinedTypeS //可以直接的隐式转换，因为不存在方法集误用的可能。
	undefinedTypeS = ms //可以直接的隐式转换，因为不存在方法集误用的可能。

	//------以下是带tag的结构体（struct）类型的转换----
	//总体来说，底层类型一样的struct类型， 如果struct类型定义不同的tag元数据，需要进行显式转换，
	//不能进行隐式转换（直接赋值），主要是为了避免对tag元数据的误操作。也就是说，tag元数据应视作
	//底层层类型的差异的参考，tag元数据不同但内存结构相同的struct允许显式转换，但不能隐式转换。
	//已定义的结构体类型,无tag原数据，底层类型是 struct { n int}
	type NoTagStruct struct {
		n int
	}

	var x struct {
		n int `foo`
	} //含有tag标记的未定义的结构体类型，如果不考虑tag元数据与y的类型相同，考虑tag元数据，则不同.
	var y struct {
		n int `bar`
	} //含有tag标记的未定义的结构体类型，如果不考虑tag元数据与y的类型相同，考虑tag元数据，则不同.
	var z NoTagStruct = NoTagStruct{n: 10}
	x = (struct {
		n int `foo`
	})(z)

	z = NoTagStruct(y)
	//x = y //编译错误，尽管二者都是未定义的struct{ n int }类型，因为tag元数据不同，不能直接转换。

	x = (struct {
		n int `foo`
	})(y) //显式的类型转换，提醒编写和代码维护的程序员，避免直接隐式转换导致对tag元数据的错误处理。

	y = (struct {
		n int `bar`
	})(x) ///显式的类型转换，提醒编写和代码维护的程序员，避免直接隐式转换导致对tag元数据的错误处理。

}

// -----------指针类型相关的类型转换
func TestUnderlyinngPointerTypeConvert(t *testing.T) {
	type MyIntType int                           //MyIntType的底层类型是int
	type YourIntType int                         //YourIntType的底层类型是int
	type IntPtrType *int                         //IntPtrType类型的底层类型为 *int，虽然内存布局相同，但二者方法集完全不同，*int指针类型所指向类型的底层类型是int
	type MyIntPtrType *MyIntType                 //MyIntPtr的底层类型为*MyIntType，虽然内存布局相同但二者方法集完全不同。*MyIntType指针类型所指向类型的底层类型是int
	var undefYip *YourIntType = new(YourIntType) //undefYip的类型是是未定义的指针类型*YourIntType，指针所指向的YourIntType类型的底层类型是int，因此可以与任何指向以int为底层类型的未定义指针类型进行显式转换。
	var pi *int = new(int)                       //pi 类型是*int，属于未定义类型(undefined type).
	var ipt IntPtrType = pi                      // ipt类型已定义类型IntPtrType，其底层类型是*int，与未定义类型*int可以相互地隐式地转换。
	var _ = ipt

	var undefMip *MyIntType = (*MyIntType)(pi) //undefMip的类型是未定义的指针类型*MyIntType，指针所指向的MyIntType类型的底层类型是int，因此可以与任何志向以int为底层类型的未定义指针类型进行显式转换。
	undefMip = (*MyIntType)(undefYip)          //两者都是未定义指针，但是各自所指向类型的底层都是int，故此，可以相互转换。
	var defMip = MyIntPtrType(undefMip)
	defMip = MyIntPtrType((*MyIntType)(undefYip)) //经过间接的显式转换，确保程序员了解到类型方法集的切换。
	_ = defMip

}

//------------------------与底层类型相关的转换规则的练习，end----------------------------------//
/******************************************************************************************
 2.3 通道（channel）类型值之间的转换规则。
     通道是一种特殊的“容器”，这个容器对其内数据元素的读、取进行了控制。因此，通道中传递的数据元素类型对
	 通道值的内存布局有影响，而通道的方向（读\取、写\发）则属于方法集。
	 因此，通道类型的转换规则是：
	 （1）通道所含元素类型不一样不能进行任何方式（隐式或显式）转换。
     （2）如果通道所含元素类型相同，但是通道的能力不同，那么，能力强的通道（双向通道）可以隐式地转换为
	     能力相同或较弱的通道（单向通道）。但是，能力弱的通道不能通过任何方式转换为能力强的通道。
*******************************************************************************************/
func TestChannelConvert(t *testing.T) {

	bdInt64Chan := make(chan int64, 10)

	var roInt64Chan <-chan int64 = bdInt64Chan // 双向(读写)channel可以隐式转换为单向(只读)channel
	//bdInt64Chan =(chan int64) roInt64Chan // 单向(只读) channel 无法转换为双向channel
	var _ = roInt64Chan
	type MyInt64 int64
	bdMyInt64Chan := make(chan MyInt64, 10)
	bdInt32Chan := make(chan int32, 10)
	var _ = bdMyInt64Chan
	var _ = bdInt32Chan
	//bdInt64Chan = (chan int64)(bdInt32Chan) //尽管int64兼容int32，但是容纳两种类型的channel不可以显式转换
	//	bdMyInt64Chan = (chan MyInt64)(bdInt64Chan)//尽管MyInt64和int64底层类型相同，但是chan int64与chan MyInt64类型值不能相互显示转换。
	//	bdInt64Chan = (chan int64)(bdMyInt64Chan)//尽管MyInt64和int64底层类型相同，但是chan int64与chan MyInt64类型值不能相互显示转换。

}
func TestUntypeConst(t *testing.T) {
	type MyString string
	const typeStr string = "hello" //这是声明了一个已定类型常量（typed const）。
	const untypeStr = "hello"      // 这是一个未定类型常量（untyped const）,可以根据需要被编译器动态转换为兼容类型。
	var ms MyString
	//ms = typeStr   虽然MyString的底层类型是string，不能把string类型的常量直接赋值给MyString，
	ms = untypeStr //可以把untype类型的string差常量赋值给string以及任何以string为底层类型的类型。
	fmt.Println(ms)

}

func TestConstantDeclaration(t *testing.T) {
	const ADMIN_ROLE_NAME = "admin"       //未定类型常量(untyped constannt)，缺省类型是string
	const DEFUALT_PSSWD string = "123456" //这是一个确定类型，string类型的常量。
	const MAX_COUNT = 120                 //MAX_COUNT是未定类型常量，缺省类型是int。根据需要，MAX_COUNT可以是int16，也可以是int32，int64，或任何以int为底层类型，能够存储值120的各种类型。
	const MIN_COUNT = 1
	const RANGE = MAX_COUNT - MIN_COUNT //这是一个untype int 类型的常量。
	//var bigFloat float64 = 1.2e10000 //变量的数值范围超边界
	const bigFloat = 1.2e10000            //1.2乘10的10000次方，该常量的数值不会引起超界异常。引入超界常量的目的主要用于如下的常量之间的计算，可保持精度
	const constRate = bigFloat / 1.2e9999 //constRate=10,两个超界常量的计算，精度更高。
}

// 未定类型常量已定类型在类型转换中的运用。
func TestConstantTypeConvert(t *testing.T) {
	const untypedIntconstant = 1       //未定类型常量 (untyped const)，缺省类型为int，所以untypedIntconstant缺省类型是int。
	const typedInt32Constant int32 = 2 //已定类型常量(typed const)，类型为int32
	var a int16 = 0
	b := a + untypedIntconstant         //常量untypedInt依从表达式中已有类型，也就是a的类型，故而b类型为int16
	fmt.Printf("var b type is %T\n", b) //打印显式类型int16
	var c int = 5
	//c += typedInt  //已定类型常量无法与类型不匹配的变量直接进行操作
	c += int(typedInt32Constant) //已定类型常量必须显式转换为目标类型才能进行操作。
	d := untypedIntconstant + 12 //  这个语句中，未定类型的常量参与表达式没有已定类型，使用缺省类型。

	fmt.Printf("var d type is %T\n", d) //变量d的类型为整数常量的缺省类型 int

	const MaxUint = ^uint(0)    //^(XOR) 运算符以单操作数出现，意味着按位取反。
	fmt.Printf("%v\n", MaxUint) //十进制打印数值
	fmt.Printf("%x\n", MaxUint) //十六进制打印数值
	fmt.Printf("%o\n", MaxUint) //八进制打印数值
	fmt.Printf("%b\n", MaxUint) //二进制打印数值
}

func TestNumbericConvertion(t *testing.T) {
	//涉及值溢出的转换。
	{
		var aInt int = 0
		var aFloat float64 = 1234.9999
		aInt = int(aFloat) //将浮点数显式转换为整数时，去掉小数部分。
		println(aInt)      //1234
		aInt = 5678
		aFloat = float64(aInt) //将整数显式转换为浮点数时，整数位保留，小数部分无限接近0.
		println(aFloat)        //5.678000e+003
		var aUint8 uint8 = 0
		var aUint16 uint16 = 0xFF01
		aUint8 = (uint8)(aUint16) //数字类型转换时允许（高位）的数值溢出，因此，aUnit8拥有低8位，故而值为1.
		println(aUint8)           //1

	}
	//涉及常量值与变量值之间的转换
	{
		const a = -1.23
		var b = a

		//var x = int32(a) //不允许将常量用于不同类型的显式转换。
		var c = int32(b) //不同类型数值变量之间则可以显式转换。
		println(b)
		println(c) //-1
	}
	//涉及常量溢出
	{
		const n = 1 << 64 //无类型常量允许值溢出缺省类型的取值范围

		//const m int = 1 << 64 //有类型常量不允许值溢出类型的取值范围
	}

}

// ------------------------下面代码是为了测试不确定类型的接口值的来源类型的精准判断-----//
type MyString string

func (ms MyString) doPrint() { fmt.Println("MyString---" + ms) }

type MyStringStr MyString

func (mss MyStringStr) doPrint() { fmt.Println("MyStringStr---" + mss) }

// func (mss MyStringStr) doNothing() {  }
type Printable interface {
	doPrint()
}

func TestRuntimeTypeConvert(t *testing.T) {
	var ms MyString = "hello"
	var s string = "world"
	var mss MyStringStr = "hello world"
	// 通过静态的强制转换，把string类型值转      换为Mystring，可行
	(MyString)(s).doPrint()
	//通过静态的强制转换，把MyString类型值转换为string，可行
	var s2 string = (string)(ms)
	_ = s2

	var ms3 MyString = MyString(mss)
	_ = ms3
	var mss2 MyStringStr = MyStringStr(ms)
	_ = mss2
	//静态可以强制转换的类型，以动态的方式无法转换。
	//msx,ok都是if语句内部的局部变量，不在函数体的范围之内。
	if msx, ok := interface{}(s).(MyString); ok {
		fmt.Println(msx)
	} else {
		fmt.Println("can not convert string to MyString dynamicly")
	}

	if sx, ok := interface{}(ms).(string); ok {
		fmt.Println(sx)
	} else {
		fmt.Println("can not convert MyString to string dynamicly")
	}
	// 不确定的 接口类型的精准类型判断
	var ptb Printable = ms //接口类型实际来源于MyString
	if mstr, ok := ptb.(MyString); ok {
		mstr.doPrint() //具体类型精准匹配，所以本语句可以执行
	} else {
		fmt.Println("can not convert Printable to MyString dynamicly")
	}
	if mssStr, ok := ptb.(MyStringStr); ok {
		mssStr.doPrint() //具体类型不匹配，所以本语句无法执行
	} else {

		fmt.Println("can not convert Printable to MyStringStr dynamicly")
	}
}

///////////////////////////////////////////////////////////////////////////
