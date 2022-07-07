package gotypeandvalue

import (
	"errors"
	"fmt"
	"strconv"
)

/*******************************************************************************************
1. 什么是类型（Type），类型有什么用途？如何定义？
类型（Types）确定了一组值（value——常量或变量）的集合，这组值的内存布局，
以及用于值集合中这些值的操作（operations）与方法（methods），也就是方法集。

可以说，“类型（Type）”是对性质相同的值或量的归纳与抽象。而“泛型(Generic Type)”则是对性质相同的类型的归纳与抽象。
有了对值的共性的归纳与抽象，就可以将一组操作或方法施加于这些值上，而不会产生错误，这就是类型的用途。
因此，在GO语言中，对类型的区分非常严格。不同类型之间严禁直接操作。

由于类型是对一组值的归纳与抽象，那么类型必须能够确定“值”在内存中格式与大小。
Go语言预定义了若干具有确定内存格式与大小的“基本类型（bool，numberic，string）”。
再以这些基本类型为基础，通过“类型组合”、以及“作为源类型（Source Type）”等两种方式衍生出新类型，从而使得衍生出的新
类型也可以确定其值的内存格式与大小。由于可以通过对包括衍生在内的所有已有类型进行“类型组合”，
或者将已有类型作为“源类型（Source Type）”的方式来衍生新类型，所以使用GO语言可以衍生出无尽的丰富类型。
类型组合方式衍生类型又称复合类型（Composite Type），Go语言定义的组合类型有8个：数组（array），结构（struct），指针（pointer），函数（function），
接口(interface),切片（slice），映射(map),通道（channel）等。 数组、切片和映射又称容器类型。
通过源类型衍生新类型方式是通过以下语法：
        type NewType SourceType
SourceType可以基本类型，基本类型的组合类型，也可以是通过 type 关键字声明的其他衍生类型。
注意：新定义的类型NewType与源类型SourceType是两个完全不同的类型，二者之前除了共享相同的底层类型外，
再无其他关系，也就是说新类型除了通过SourceType推导出底层类型，从而决定所属类型值的内存布局之外，
与SourceType再没有任何瓜葛，二者之前相互不共享、不继承“操作”与“方法”的集合。

后续会详细介绍什么是底层类型，以及如何推断底层类型。

***********************************************************************************/
//go 类型的严格性体现
func TestTypeOperation() {
	var a int = 12
	var b int16 = 24
	//c := a + b //编译不过去，GO认为C的类型无论设置为int或int16都容易引起误解。
	c := int(b) + a //必须显式地将b转换为与a一样的类型才能相加。
	c++
}

//----------------------------Type Define Exercise begin ----------------------------//
type MyFloat float64                      //float64作为MyFloat的源类型。
type NewFloat MyFloat                     //MyFloat作为NewFloat的源类型。
var intArrayValue [6]int                  //通过组合类型（也就是字面类型）[6]int 定义一个了数组变量。
type IntArray [6]int                      //这里以组合类型（也就是字面类型）[6]int  作为源类型定义了一个新类型。
var intSliceValue []int64                 //这里以组合类型（也就是字面类型）[]int64（int64切片）定义了一个变量。
type Int64Slice []int64                   //这里以组合类型（也就是字面类型）[]int64（int64切片作为源类型定义了一个新类型。
var floatPointerValue *float64            //通过组合类型（也就是字面类型）*float64  定义一个了数组变量。
type FPointer *float64                    //通过类型组合产生的字面类型（无名字类型）*float64作为源类型定义了一个新类型的命名类型FPointer，  字面类型*float64 是 * 与 float64的组合 。
var myIntToStrFunc func(int) string = nil //直接使用组合类型（没有名字的字面类型）func(int) string 定义变量，注意不是定义类型。
type IntToStrFunc func(int) string        //通过类型组合方式产生的字面类型
var studentData struct {
	id   string
	name string
} // 这里直接使用组合类型 struct{id  string  name string} 定一个变量,注意，不是定义了类型。

type Student struct {
	id   string
	name string
} //这里以组合类型（也就是字面类型）struct{id  string  name string}作为 源类型，定义了一个新的命名类型Student。

var studentMapData map[string]Student //这里以组合类型（也就是字面类型）map[string]Student定义了一个变量。
type StudentMap map[string]Student    //这里以组合类型（也就是字面类型）map[string]Student作为源类型定义了一个新类型。
var readerObj interface {
	read([]byte) (int, error)
} = nil //这里用组合类型（也就是字面类型） interface { read([]byte) (int, error)} 定义了一个变量。

type Reader interface {
	read([]byte) (int, error)
}                              //这里用组合类型（也就是字面类型） interface { read([]byte) (int, error)} 作为源类型，定义了一个新类型。
var intDataReadChan <-chan int //这里以组合类型（也就是字面类型）<-chan int 定义了一个变量。
type IntReadChan <-chan int    //这里以组合类型（也就是字面类型）<-chan int 作为源类型定义了一个新类型。

type IntList = []int //注意，这里为组合类型（也就是“字面类型” ）[]int定义了一个别名。IntList只是组合类型（也就是“字面类型” ）[]int的别名，本质还是组合类型。

//----------------------------Type Define Exercise end  ----------------------------//

/**************************************************************************************
2.类型如何分类？

无论是类型还是值，只要需要反复多次引用，就需要使用“标识符”作为名字来定义这个类型或值。如果无需反复多次引用，
那么就不需要命名。从这个角度来看，类型分为命名的类型(Named Type)和未命名类型(UnNamed Type)。
有名字的类型是用type NewType SourceType来定义的类型。所有预定义的类型都是命名类型（Named Type）。
而所有组合类型，也就是字面类型都是未命名类型（UnNamed Type）。
在GO1.9之前就是通过是否有名名字来对类型进行分类，但是随着可以为类型定义“别名（alias）”之后，按照是否有名字给
类型分类就会产生语义上的混淆。因为使用 Type Alias = T定义一个名字标识符时，如果T是未命名类型（类型字面量）时，
尽管Alias像是是“类型名”，但所代表的类型仍然是未命名类型，只有当T是命名类型时，Alias才代表一个命名类型。
所以，在GO1.9之后，使用“已定义类型(Defined Type)”和“未经定义类型（Undefined）”来对类型分类，
“已定义类型(Defined Type)”：是指通过“Type NewType SourceType” 语句明确定义过的类型，以及GO语言预定义类型，以及这些类型的别名。
“未经定义类型”：是指用到时才书写的字面类型，及其别名。
其实，我们不要将type alias=T 语法所定义的“类型的别名（alias）”看作一个类型，它只是所表示“类型”的另一种写法。
这样，已定义类型就等同于了命名类型(name type)，而未经定义类型就等同于了未命名类型（unnamed type）.
换句话说，类型别名所表示的类型是命名类型（已定义类型）还是未命名类型（未经定义类型）取决该类型的真实定义情况。
当我们谈论一个类型的名称时，这个名称可能是 一个已定义类型的名称，也可以是一个别名。

*********************************************************************************************/

//---------------------------- Type Classify Exercise begin ----------------------------//
type A []string   //A是已定义类型，即：go1.9之前的命名类型（名字为A），[]string是一个“未经定义类型”，或“字面类型”，即：go1.9之前的未命名类型。
type B = A        //B是类型A的别名（类型A的本命是A，现在有了别名B），由于A是已定义类型，类型名B实际上表示的是类型A，所以类型名B所代表的类型是已定义类型（命名类型）
type C = []string //C是"未经定义类型"，或者“字面类型”[]string的别名，因此，C所代表的类型是“未经定义类型”或“未命名类型”

//下面语句用来声明各种已定义类型（GO1.9之前的命名类型）
type MyStringType string
type Person struct { //这里使用了未命名类型，也就是“字面类型”struct {Id int64,Name string}作为命名类型 Person的底层类型。
	Id   int64
	Name string
}
type Printer interface { //这里使用了未命名类型，也就是“字面类型”nterface {...}作为命名类型Printer的底层类型。
	Print(data interface{}) (bool, error)
}
type StringSlice []string //这里使用了未命名类型，也就是“字面类型” []string作为命名类型StringSlice的底层类型。
func ShowUnderlyingType() {
	var p = Person{Id: 1, Name: "Lantian"}

	fmt.Printf("Person类型的实际类型为%T\n", p)

}

//下面使用“未经定义类型”，即，"字面类型(Type literal)"，也就是Go1.9之前的“未命名类型(Unnamed Type)”，来定义变量。
var unDefinedTypeVar1 []string //这是一种常见的方式，使用未命名类型，也就是类型字面量定义了一个字符串切片类型的变量。

var unDefinedTypeVar2 struct { // 这是一种不常见变量定义方式，往往以struct{...} 为底层类型的类型都是命名类型。
	Id   int64
	Name string
}
var unDefinedTypeVar22 = struct { // 这是一种不常见的结构体变量定义并赋初值的方式，往往以struct{...} 为底层类型的类型都是命名类型。
	Id   int64
	Name string
}{
	Id:   5,
	Name: "刘邦",
}

var unDefinedTypeVar3 interface { //这是一种不常见的方式，往往以interface{...} 为底层类型的类型都是命名类型。
	Print(data interface{}) (bool, error)
}

var unDefinedTypeVar4 func(int) (int, error) //这是一种不常见的变量定义方式。

func UseUnnamedTypeVar() {
	var person1 = Person{Id: 1, Name: "lantian"}
	var person2 Person = unDefinedTypeVar22
	println("person1=Person { Id:" + strconv.FormatInt(person1.Id, 10) + " , Name:" + person1.Name + " }")
	println("person2=Person { Id:" + strconv.FormatInt(person2.Id, 10) + " , Name:" + person2.Name + " }")
	unDefinedTypeVar2 = person1

	println("unnameTypeVar1= { Id:" + strconv.FormatInt(unDefinedTypeVar2.Id, 10) + " , Name:" + unDefinedTypeVar2.Name + " }")
	unDefinedTypeVar2.Name = "liufei"
	println("变更后，unnameTypeVar1= { Id:" + strconv.FormatInt(unDefinedTypeVar2.Id, 10) + " , Name:" + unDefinedTypeVar2.Name + " }")
	fmt.Printf("person1的类型是：%T \n", person1)
	fmt.Printf("person2的类型是：%T \n", person2)
	fmt.Printf("unnameTypeVar2的类型是：%T\n", unDefinedTypeVar2)
	fmt.Printf("unnameTypeVar22的类型是：%T\n", unDefinedTypeVar22)
	unDefinedTypeVar4 = func(i int) (int, error) {
		if i > 100 {
			return 0, errors.New("bad input")
		} else {
			return i * 100, nil
		}
	}
	num, _ := unDefinedTypeVar4(20)
	println("用参数20调用函数变量unnnameTypeVar4，调用结果是：", num)
	fmt.Printf("unnnameTypeVar4的类型是：%T \n", unDefinedTypeVar4)

}

//------------------------------Type Classify Exercise begin end -------------------------------//

/***********************************************************************************************
  三、什么是类型的方法集（Method sets），如何定义，有什么用途？
	每个类型都有一个方法集（可能为空）与之相关。
	类型方法集的定义：必须是与已定义类型（defined type）在同一个包中（方便编译器搜寻类型的方法集），以该类型
	作为“接收者”参数的函数就是类型的方法集。
	注意，未经定义的类型（类型字面量）或等价的别名类型都不能定义方法集。因为类型字面量可以在任何包中随处写，
	编译器无法确定其方法集到底包含多少个函数，故而只允许为已定义类型定义方法集。
	在类型的方法集中，每个方法都必须有一个唯一的“非空格（non-blank）”的方法名。

	类型的方法集的用途主要是用于决定类型所实现的接口（用于多态），并且方法集中的方法可以使用该类型的接受者实例进行调用。


  方法集范围判定规则：
	1.接口（interface）类型的方法集就是接口自己。
	2.对于任何一个非接口类型的其他类型T，其方法集范围包括所有声明为“以类型T为接收器”的方法。
	3.而类型T所对应的指针类型 *T（注意，*T是不同于T的一个新类型），其方法集包括所有声明
	为“以类型 *T和T为接收器”的方法，也就是说*T类型的方法集包括了T的方法集。
	4.对于包含了嵌入字段（embedded fields）的struct类型的更多规则会在struct类型规范中有具体介绍。
	5.除了上述情形以外的其他情况下，类型的方法集就是空的。

   在 Go 语言中，使用“已定义类型”，也就是type TA TB 语句所定义的类型有两个主要原因，
   即：（1）防止意外组合两种不同类型的值，以及（2）区分不同类型上的方法集定义。

   所以，当我们写下type A B 语句来定义新类型A时，类型A没有类型B的任何方法。
   当然，需要显式转换才能在两种类型的值之间进行转换（无论它们是否为接口类型）。
   在上面提到的两种情况（类型转换与方法集）中，A 和 B 之间都没有任何特殊关系（只是内存存储结构相同）。
*********************************************************************************************************/
//-----------------------Method Set Exercise begin------------------------//
type TypeB = []int //TypeB 等价于[]int,实际上是一个未经定义的字面类型，不允许定义方法集。
/** TypeB定义不允许定义方法集。
func (b TypeB) fb(index int) int {
	return b[index]
}
**/
type T string
type Hello interface {
	SayHello()
}

type T2 T //T2以T为源类型（Source Type），但这不意味着T2继承了T的方法集，某个类型的方法集与其引用类型的方法集之间毫无关系。

// 实现了error接口
func (t *T) Error() string {
	return string(*t)
}
func (t T) SayHello() {
	println("hello world")
}
func ShowTypeMehodSet() {
	var t T = "value of type T" // 字面量"type"的缺省类型是string，可以赋值给所有以string为底层类型的变量。
	var pt *T = &t

	//t2.SayHello()  //
	// 类型T的值t 可以调用类型 *T的方法集中（含T的方法集）中的方法，但这只是“语法糖”，并不表示类型T与类型*T的方法集相同。

	println("t.Error() call  : ", t.Error()) //t调用了*T的方法（语法糖）。
	t.SayHello()                             //t调用了类型T的方法。
	println("*t.Error() call  : ", pt.Error())
	pt.SayHello() //*T类型的pt变量调用自己的子方法集中的方法，该子方法集来自于类型T。
	testMehodSet := func(i int) (Hello, error) {

		if i <= 0 {
			//return pt, t //T类型的方法集只有SayHello()方法，只实现了Hello接口，未实现error接口。
			return pt, pt //*T类型方法集包括了Error()和SayHello()方法，因此，实现了Hello与error两个接口。
		} else {
			return t, nil
		}
	}
	testMehodSet(0)
	var t2 T2 = "value of Type T2"
	println(t2)
}

//-----------------------Method Set Exercise End------------------------
/*********************************************************************************************
  四、何为类型的底层类型？知道类型的底层类型有什么用途？

   在GO中，一个类型的底层类型（Underlying Type）决定了该类型的值在内存中的实际结构。

   因此，类型的底层类型对于于值（变量、常量、字面量）之间的转换（Converson）、
   赋值（Assignment ）与比较（Comparison ）非常重要，关于值的转换、赋值和比较
   在govalues.go文件中有详细介绍。

   每个类型T都有且只有一个底层类型，如果T是预先定义的boolean、 numeric、string 类型,
   或者“字面类型”以及unsafe.Pointer类型，那么类型T的底层类型就是它自己。
   否则，T的底层类型就是在其类型声明中的源类型（Source Type）的底层类型。
   这就需要沿着类型声明的“源类型链”来推导底层类型，推导的规则是：

   沿着类型的“源类型链”上遍历类型，如果遇到的类型是预先定义的“基本类型”，或者“未经定义类型（undefine type）”，
   也就是“组合类型”（比如[]T,Map[A]B），就终止遍历，此时得到的类型就是该类型的底层类型。注意，遇到的“组合类型”
   未毕是预先定义的基本类型的组合，遇到的是衍生类型的组合类型也要终止遍历 （因为，组合类型可以根据被组合的类型的底层类型来确定存储结构）。
   因此，类型的底层类型都是预先定义的boolean、 numeric、string类型或“字面类型（组合类型）”。
   由于一个类型与其源类型值集合，以及操作和方法集合 完全没有任何关系，二者之间只是共享了同一个底层类型，
   所以，使用 type A B 这样的写法来定义类型A，与使用type A <B 的底层类型> 这样的写法来定义类型A的方式完全相同。
**/
//--------------------------Underlying Type Exeercise begin---------------------------//
// 以下类型的底层类型都是int.
type (
	MyIntType int
	Age       MyIntType
)

//以下类型有不同的底层类型。
type (
	IntSlice   []int       // 底层类型是 []int
	MyIntSlice []MyIntType // 底层类型是 []MyInt ，遇到的引用类型是组合类型而截止。
	AgeSlice   []Age       // 底层类型是 []Age， 遇到的引用类型是组合类型而截止。
	Ages       AgeSlice    //底层类型是[]Age ，遇到的引用类型是组合类型而截止。
)

//上述类型的底层类型推导过程如下，#号表示不再继续向上推导
//		MyInt → int
//		Age → MyInt → int
//		IntSlice → []int
//		MyIntSlice → []MyInt → #[]int
//		AgeSlice → []Age → #[]MyInt → #[]int
//		Ages → AgeSlice → []Age → #[]MyInt → #[]int
//
//--------------------------Underlying Type Exeercise end---------------------------//
