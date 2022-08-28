/**
1.什么是接口类型？

  接口类型定义了一个方法\函数原型（prototype）的集合，换句话说，接口通过罗列函数原型，定义了一个方法集，
  因此，我们可以把接口类型看作一个方法集合。从另外一个角度，也可以将接口类型看作一个行为(behavior)集合。

2.接口类型的自身结构，即：接口类型值的内存结构。
 // 可能的实现方式1：
type empty_interface struct {
	dynamicType  *_type         // 描述具体类型的元数据（含方法集）。
	dynamicValue unsafe.Pointer // 实际值的内存地址（直接值部的地址）
}
//可能实现方式2：
type non_empty_interface struct {
	dynamicTypeInfo *struct {
		dynamicType *_type       //  描述具体类型的元数据（除方法集之外的部分）
		methods     []*_function //  描述接口方法集合，只有方法集合匹配的具体类型才能赋值给接口类型。
	}
	dynamicValue unsafe.Pointer // // 实际值的内存地址（直接值部的地址）
}

接口（interface）类型，在Go语言中具有若干重要作用.但是，其首要作用就是接口（interface）类型
使得GO语言支持对值（value）的“装箱（boxing）”。正是这种对值（value）的“装箱（boxing）”功能，
从而使GO语言支持多态和反射。

3.接口类型的字面表示：
   interface{
	   methodSignature1
	   methodSignature2
	   .....
	   methodSignatureN
   }
  注意 : interface{} 表示该接口类型的方法集为”空“。”空方法集“是任意一个方法集的真子集。所以，任何具体类型
  的实例都可以赋值给interface{} 类型的变量。

4.在运行时判断接口类型值的具体类型。
      ConcreteTypeValue,boolValue=InterfaceValue.(ConcreteType)
	这个语法的语义是在运行期，判断不确定具体类型的接口值InterfaceValue其来源的真实类型是否为指定的
	具体类型ConcreteType，如果是，就将该接口值转换为ConcreteType，并返回 (ConcreteTypeValue,true)，
	否则返回（零值，false）。这是一种无法完全保证类型安全的动态类型断言，因为必须依靠程序员检查转换结果。
**/

package interfaces

import (
	"fmt"
	"testing"
)

type I interface {
	M() string
}
type MyString string

// 此非正式程序中，为了避免空指针调用抛出panic异常，
// 用空指针调用时返回"NIL"字符串。
func (ms *MyString) M() string {
	if ms == nil {
		return "NIL"
	} else {
		return string(*ms)
	}
}

// fmt格式化可以打印出变量的值与类型。
func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}
func TestUseInterfaceValue(t *testing.T) {
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

func TestUseEmptyInterface(t *testing.T) {
	var i1, i2 interface{}
	i1 = "hello"
	i2 = 100
	describe(i1)
	describe(i2)
	switch t := i1.(type) {
	case string:

		println("变量i1的底层类型为string,值的内容为%v", i1)
		describe(t)
		//此时，变量t的类型为string。
		var s string = t
		println(s)
	default:
		println("变量i1的底层类型值未知")
	}
}

//-------------------------- 接口值的比较-----------------------------------//
/***
    知识点，两个接口值之间可以进行比较操作（==，!=）。
	相等的要求两个接口值的直接值部的内存内容相同。
	即：两个接口值所封装的具体类型相同，并且，指向实际值的指针也相同。
***/
type MyInt int

func (mi MyInt) doPrint() { fmt.Printf("MyInt:%v\n", mi) }

type YourInt int

func (mi YourInt) doPrint() { fmt.Printf("YourInt:%v\n", mi) }

type Printable interface {
	doPrint()
}

func TestInterfaceCommpare(t *testing.T) {

	var prt1, prt2, prt3, prt4 Printable
	var myInt1 MyInt = 100
	var myInt2 MyInt = 100
	var myInt3 MyInt = 200
	var yourInt1 YourInt = 100
	prt1 = myInt1
	prt2 = myInt2
	prt3 = myInt3
	prt4 = yourInt1
	fmt.Printf("prt1==prt2 ? : %v\n", prt1 == prt2) //true，接口值的具体类型与实际值指针都相同
	fmt.Printf("prt1==prt3 ? : %v\n", prt1 == prt3) //false,接口值的实际值指针不同
	fmt.Printf("prt1==prt4 ? : %v\n", prt1 == prt4) //false, 接口值的具体类型不同
	var ifv1, ifv2, ifv3 interface{}
	ifv1 = 100
	ifv2 = 100
	ifv3 = "hello"
	fmt.Printf("ifv1==ifv2 ? : %v\n", ifv1 == ifv2) //true，接口值的具体类型与实际值指针值都相同
	fmt.Printf("ifv1==ifv3 ? : %v\n", ifv1 == ifv3) //false ,接口值的具体类型不同
}
