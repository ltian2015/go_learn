/**
关于接口值的知识点：
知识点1:接口也是值。它们可以像其它值一样传递。接口值可以用作函数的参数或返回值。
在GO语言内部，接口值可以看做包含值和具体类型的元组：
(value, type)
接口值保存了一个具体底层类型的具体值。接口值调用方法时会执行其底层类型的同名方法。
底层值为 nil 的接口值，即便接口内的具体值为 nil，方法仍然会被 nil 接收者调用。
在一些语言中，这会触发一个空指针异常，但在 Go 中通常会写一些方法来优雅地处理它（如本例中的 M 方法）。
注意: 保存了 nil 具体值的接口其自身并不为 nil。
知识点2: 指定了零个方法的接口值被称为 “空接口”类型，定义如下：
interface{}
空接口这种类型的值可保存任何类型的值。（因为每个类型都至少实现了零个方法。）
空接口类型相当于java语言中的Object类型，因而空接口被用来处理未知类型的值。
例如，fmt.Print 可接受类型为 interface{} 的任意数量的参数。
**/
package interfaces

import "fmt"

type I interface {
	M() string
}
type MyString string

//此非正式程序中，为了避免空指针调用抛出panic异常，
//用空指针调用时返回"NIL"字符串。
func (ms *MyString) M() string {
	if ms == nil {
		return "NIL"
	} else {
		return string(*ms)
	}
}

//fmt格式化可以打印出变量的值与类型。
func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}
func UseInterfaceValue() {
	var ms *MyString
	var i I = ms //经过这样的赋值后，i的值为二元组(nil,*interfaces.MyString)
	var i2 I     //i的值为nil
	//验证i值是否为nil。
	if i == nil {
		println("接口值i为nil") //不会执行，因为i值不是nil，尽管i值二元组中的底层值是nil
		return
	}
	fmt.Printf("(%v, %T)\n", i, i)
	println(i.M())
	//验证i2是否为nil
	if i2 == nil {
		println("接口值i2为nil") //会执行，因为i2确实为nil
		return
	}
	fmt.Printf("(%v, %T)\n", i2, i2)
}

func UseEmptyInterface() {
	var i1, i2 interface{}
	i1 = "hello"
	i2 = 100
	describe(i1)
	describe(i2)
	switch t := i1.(type) { //var.(type),只能用在case语句中的类型断言
	case string:

		println("变量i1的底层类型为string,值的内容为%v", i1)
		describe(t)
		//此时，变量t的类型为string。
		var s string = t
		println(s)
	default:
		println("变量i1的底层类型值未知")
	}
	var i = i2.(int) //类型断言,注意与i1.(type)的区别，v.(type)语句只能用在switch语句
	println(i)
}
