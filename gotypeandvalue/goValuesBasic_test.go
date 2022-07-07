package gotypeandvalue

import (
	"errors"
	"fmt"
	"testing"
	"unsafe"
)

/*************************************************************************************************************
1.有关值（Value）的基本知识。
（1）值的表达方式：在go语言中，值（Value）有字面量（无名常量）、命名常量、变量和表达式四种，前三者都可以看作“表达式”的特殊形式。

（2）不同类型的字面量：如果说类型存在“已定义类型(defined type =named type)”和未经定义类型(undefined type =unnamed type),
值也分为有名字的值与无名字的值。为值设定名字（标识符）是为了反复多次引用，有名字的值又分为“常量”和“变量”，
无需被多次引用的值或量可以直接书写出来，无需命名，被称为“字面量（Value literals）——无名常量”。
系统预先定义的类型bool，string,numberic都支持字面量，而函数（function）、容器类型，也就是
数组（slice）、切片（Slice）、Map与结构体（struct）等组合类型也都可以写出字面量。但是，指针（Pointer）、
通道（channel）、接口（interface）无法写出字面量。因为指针类型的变量存储的数据是运行时的内存地址，
程序员在编写程序时无法给出。而通道（channel）则是一个数据通信的协调器，不是一个数据，无需给出字面量。
接口（interface）则表示的是具有共性操作或方法集的一个抽象概念，它总是由具象类型所实现，因此，也无法给出字面量。

（3）字面量的类型。在计算机中，如果不知道值的类型就无法对其进行操作，因为不知道值的内存结构的布局以及适用该值方法或操作集合。
但是在GO语言中，有一类值是没有类型的（untyped value），或者说其类型是动态的。因为，不确定类型反而有利于内存优化。
这些值就是“字面量”或没有定义类型的命名常量。因为，字面量或常量的内存大小是在编译期确定，而且一旦内存确定就永远不变化。
比如，无类型的常量 5（假定， const i=5），在编译时只需要用一个字节存储即可，但是一旦指定其类型为int64，就要分配8个字节的内存来存储。
因此，不确定类型非常有利于内存优化，这样，在使用改常量的时候即可将其当作int8，也可当作int16，还可以当作int32或int64。
但是这些字面量一般都有一个缺省类型，比如，再需要推断类型时，就是用它们的缺省类型。而变量由于存储的值会变化，必须指定内存的边界，也就是
内存的格式或布局，所以，变量必须都有类型，不能像未定类型的常量那样存在未定类型的变量。
但是有一个特殊的字面量除外，字面量 nil 的类型是动态的，可以是funtion，也可以map，但是字面量nil没有缺省类型。这些没有定义类型或者
动态类型的值在使用时按照缺省类型对待。对于明确了类型的值，我们称为类型明确值（typed value），类型明确值主要是变量，和声明了类型的常量，
以及含有类型明确值的表达式。
以下是各种未定类型常量的缺省类型：
1.字符串常量，比如：字面量 "hello" 的缺省类型是string
2.整数常量，比如，字面量  12 的缺省类型是int
3.浮点数常量，比如字面量  1.2的缺省类型是float64
4.rune（int32的别名），比如字符字面量  '我' 的缺省类型是int32
5.复数虚部常量，比如： 字面量 0i的缺省类型是complex128

（4）类型的零值。每个类型有一个零值。一个类型的零值可以看作是此类型的默认值。 预声明的标识符nil可以看作是切片、映射、
函数、通道、指针（包括非类型安全指针）和接口类型的零值的字面量表示。而预定义的类型的零值，比如string，numeric，bool的
零值则不是nil，分别是 "",0,false。前面说过，nil是个特殊的字面量，它没有缺省类型。
(5). 为了内存安全，常量、字面值以及除了变量之外的表达式，都不支持取地址操作 & 。

************************************************************************************************/
//------------------值的基本知识练习，开始-----------------------------------------------------//
//未定类型常量、字面量的定义与使用
const MAX_AGE = 300 //未定类型常量300,也就是动态的整数类型的常量，该常量占用2个字节空间，可以赋值给任意内存空间兼容的整数类型。
//var m int8 = MAX_AGE // 编译错误，类型的内存空间不兼容，int8变量只有一个字节，无法容纳两个字节常量值。
var n int16 = MAX_AGE //
var o int32 = MAX_AGE
var p int64 = MAX_AGE

//var pi *int = &MAX_AGE //常量不支持取地址操作。

//var pi *int64 = &(p - int64(0)) //除了变量意外的表达式不能取地址
var INT_AGE = 15 //尽管INT_AGE没有显式地指定类型，但是变量类型是int，因为 未定类型字面量15的缺省类型是int
//var i int64 = INT_AGE //编译错误，类型不匹配，未经转换不能直接赋值。
var j int = INT_AGE

const MAX_SPACE = 65535 //占满了两个字节的未定类型常量。
//var oo int16 = MAX_SPACE + MAX_AGE //编译错误，常量表达式在编译期已求值，结果占用空间大于2个字节。

//打印出各种未定类型常量的缺省类型。
func TestUntypedConstantDefaultType(t *testing.T) {
	fmt.Printf("default type of constant %v is %T\n", "hello", "hello") // int类型
	fmt.Printf("default type of constant %v is %T\n", 0, 0)             // int类型
	fmt.Printf("default type of constant %v is %T\n", 0.0, 0.0)         //float64类型
	fmt.Printf("default type of constant %v is %T\n", 'x', 'x')         //int32 类型
	fmt.Printf("default type of constant %v is %T\n", 0i, 0i)           //complex128
	fmt.Println("--------------------------------------------------------------------")
	s := "hello"
	a := 0
	b := 0.0
	c := 'x'
	d := 0i
	fmt.Printf("var  s type is  %T,get default type from constant %v\n", s, "hello")
	fmt.Printf("var  a type is  %T,get default type from constant %v\n", a, 0)
	fmt.Printf("var  b type is  %T,get default type from constant %v\n", b, 0.0)
	fmt.Printf("var  c type is  %T,get default type from constant %v\n", c, 'x')
	fmt.Printf("var  d type is  %T,get default type from constant %v\n", d, 0i)
}

//------------------------------------值的基本知识练习，结束----------------------------------//

/************************************************************************************
2.GO语言常量的类型限制。
golang中const只支持预定义的基本类型的值（numberic系列类型，string类型，bool类型）,
或以基本类型为底层类型的命名类型。
不支持未经定义类型（非基本类型的复合类型复合类型，字面类型）的常量，也不支持以未经定义类型（非基本类型的复合类型）作为源的
已定义类型常量。
***********************************************************************************/
//-----------------------GO语言常量的类型限制 练习 开始-------------------------------------//
const NAME1 string = "KATE"    //ok，基本类型
type NameType string           //以基本类型为源的已定义类型（defined type）
const NAME2 NameType = "KATE2" //ok,支持以基本类型为源的“定义化类型（defined type）”
type StudentType struct {
	id   string
	name string
}

/**
const STUDENT1 StudentType = StudentType{
	id:   "01",
	name: "KATE",
}
**/
//非基本类型或非以基本为源的定义化类型无法定义常量，所以上面的STUDENT1不能定义为常量，只能定义为如下的变量。
var STUDENT1 = StudentType{id: "001", name: "liufei"}

//下面的EOF变量模仿了io包中的io.EOF变量。EOF本应是常量，
//但由于常量不支持非基本类型，所以接口类型的error不能直接定义为常量。
var EOF error = errors.New("End Of File")

//为了让error可以作为常量，下面给出了比较好的error常量定义的方法,
//MyError以string为底层类型，实现了error接口。
type MyError string

func (me MyError) Error() string {
	return string(me)
}

const MY_EOF MyError = "EOF-End Of File"

func TestMyError(t *testing.T) {

	errCreate := func(i int) error {
		if i <= 0 {
			return MY_EOF
		} else {
			return nil
		}

	}
	err := errCreate(0)
	if err != nil {
		fmt.Println(err.Error())
	}
}

//-----------------------GO语言常量的类型限制 练习 结束-------------------------------------//

/********************************************************************************************
3.unsafe包中函数求值的常量特征。unsafe包中的函数是在编译期间求值，因此，可以作为常量表达式的组成部分。
  只要unsfe包中的函数有返回值，且是可以定义为常量的基本类型。
**********************************************************************************************/
//---------------unsafe包函数求值的常量特征练习 开始--------------------------------------------------//
//const INT_SIZE uint = getIntSize() //运行期求值的函数不能作为常量表达式的一部分。
const FLOAT_SIZE uint = uint(unsafe.Sizeof(0.11)) //unsafe包函数在编译期求值，可以作为常量表达式的一部分。

func getFloatSize() uint {
	var i float64 = 0.11          //变量都是确定类型的值，因此内存布局是编译器可知的。
	return uint(unsafe.Sizeof(i)) //由于i的内存布局可知，尽管i是变量，在编译期间，unsafe.Sizeof(i)仍能求值。而函数整体在运行期相当于返回一个常数。
}

//---------------unsafe包函数求值的常量特征 end--------------------------------------------------//

/***************************************************************************************************
4、iota【a:juta】在常量声明中的用途。
 iota是第9个希腊字符的（大写Ι，小写ι）的名字。通常用于表示数学计数。
 iota是builtin.go中定义的系统内置未定义类型常量，值为0 .
 含义来源可见https://stackoverflow.com/questions/31650192/whats-the-full-name-for-iota-in-golang
 在golang中，iota主要用于简化“批量式”的常量定义，iota出现在一个有效常量定义行中的时候，它表示计数0，
 后面的有效常量定义行（非注释和空行）不用写代码，编译器会自动复制上一行中定义常量的含iota常量的常量表达式，
 并将iota计数加1。当然，如果常量不在同一批中，则iota互不相让影响。
**********************************************************************************************/
//--------------------------------- iota 在定义常量中的用途示例 开始---------------------------//
func TestConstDefineUsingIota(t *testing.T) {

	const ( //这是声明的第一批常量，共6个有名称的常量和2个被空标识符定义的名称待定或作废常量
		ONE   = 1 + iota //1+0
		TWO              // 编译器自动复制上一个显式的iota常量表达式，此时iota=1使得TWO=1+1
		_                // 编译器自动复制上一个显式的iota常量表达式，此时iota=2使得_=1+2
		_                // 编译器自动复制上一个显式的iota常量表达式，此时iota=3使得_=1+3
		FIVE             //编译器自动复制上一个显式的iota常量表达式，此时iota=4使得FIVE=1+3
		SIX   = 1 + iota //此时，iota为5，故SIX=1+5=6
		EIGHT = 2 + iota //此时，iota为6，故EIGHT=2+6=8
		NINE             //编译器自动复制上一个显式的iota常量表达式，此时iota=7使得NINE=2+7
	)
	const ZERO = iota //声明的第二批常量，共1个
	const (           //声明的第三批常量，共3个有名称的常量和1个空标识定义的名称待定或作废常量。
		_  = iota             //0
		KB = 1 << (10 * iota) //左移n位相当于2的n次方，右移n位相当于除以2的n次方，本表达式相当于2的10*1次方，
		MB                    //编译器自动复制上一个显式的iota常量表达式，此时iota=2，MB =1 << (10 * 2)，即：2的20次方。
		GB                    //编译器自动复制上一个显式的iota常量表达式，此时iota=3，MB =1 << (10 * 3)，即：2的30次方。

	)
	const NAME_KITE string = "kate"

	fmt.Printf("ZERO=%v\n", ZERO)
	fmt.Printf("ONE=%v\n", ONE)
	fmt.Printf("TWO=%v\n", TWO)
	fmt.Printf("FIVE=%v\n", FIVE)
	fmt.Printf("SIX=%v\n", SIX)
	fmt.Printf("EIGHT=%v\n", EIGHT)
	fmt.Printf("NINE=%v\n", NINE)
	fmt.Printf("KB=%v\n", KB)
	fmt.Printf("MB=%v\n", MB)
	fmt.Printf("GB=%v\n", GB)
}

//--------------------------------- iota 在定义常量中的用途示例 开始---------------------------//
/*******************************************************************************

***********************************************************************************/
