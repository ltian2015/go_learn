package generic

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
)

// 泛型中，无论是泛型函数，还是泛型类型，都是定义类型参数列表。
// 定义类型参数与定义普通参数一样，都要定义参数名及对应的参数类型。
// 类型参数的被称为类型是元类型(meta-type)，也就是关于类型的类型。
// 元类型在go泛型中被称为约束“constain”
// 泛型类型方法(method)的接收者（reciever）的类型应与该泛型类型一样。
// 但方法接收者的泛型类型中的类型参数名称可以与泛型类型定义中的类型参数名称不同。
// 比如，泛型类型Range的类型参数名字为E，该类型参数E的元类型（meta-type）是ordered，
// 其方法Handle中接收者参数的中的类型参数名称为T，也就是说接受者的类型是Range[T],
// 但这个类型参数T的元类型（meta-type）仍然是ordered。说到底参数类型名称只是为元类型（meta-type）
// 提供了一个便利的引用名而已，元类型（meta-type）才是重要的。
// 在使用泛型的时候，一定要为泛型的形式参数（parameter）指定具体类型的实际参数（argument）。
// 这种操作被称为泛型类型的实例化（instantiation）。泛型类型（Generic）实例化的结果是产生了一个
// 类型（Type）,类型(Type)实例化后，会产生值(Value)。
//
// ListHead is the head of a linked list.
type ListHead[T any] struct {
	head *ListElement[T]
}

// ListElement is an element in a linked list with a head.
// Each element points back to the head.
type ListElement[T any] struct {
	next *ListElement[T]
	val  T
	// Using ListHead[T] here is OK.
	// ListHead[T] refers to ListElement[T] refers to ListHead[T].
	// Using ListHead[int] would not be OK, as ListHead[T]
	// would have an indirect reference to ListHead[int].
	head *ListHead[T]
}

/**
The elements of an ordinary interface type are method signatures and
embedded interface types.
 We propose permitting three additional elements that may be used in an interface type
 used as a constraint. If any of these additional elements are used,
 the interface type may not be used as an ordinary type, but may
  only be used as a constraint.
**/

func Min[T interface {
	int | int8 | int16 | int32 | int64 | float32 | float64 |
		uint | uint8 | uint16 | uint32 | uint64 | string
}](ts []T) T {
	var t = ts[0]
	var i int = 8
	// var o ordered = i // 通过三种特殊方式定义的接口，只能作为泛型约束的接口不能向普通接口那样使用。
	//_=i
	_ = i
	return t
}

// 嵌入了约束型的接口就只能作为约束而存在，不能当作普通接口来使用。
// 不能用来声明变量
type ComparableAdder interface {
	comparable //这是一个约束类型的接口，非普通接口。
	Add() int
}

func compareThenAdd[T ComparableAdder](t1, t2 T) int {
	if t1 == t2 {
		return t2.Add()
	} else {
		return t1.Add()
	}
}

////////////////////////////类型推断/////////////////////////////////////////
///////////////基于函数调用的实际参数的类型推断//////////////////////////////////////////
//对于泛型函数而言，只有当泛型类型用于函数输入参数时，在函数调用时才可以进行类型推断。
//泛型类型仅用于函数的结果或者函数体内时，在函数调用时无法推断类型。如下函数将泛型类型用于函数返回值。

func convertTo[T interface {
	int8 | int16 | int32 | int | int64
}](i int8) T {
	return T(i)
}
func Map[E, T any](es []E, f func(e E) T) []T {
	var result []T
	for _, e := range es {
		t := f(e)
		result = append(result, t)
	}
	return result
}

func TestTypeInferenceByFuncCall(t *testing.T) {
	var i8 int8 = 8
	//  因为泛型类型用于仅用于函数的返回值，调用时必须显式地给出类型参数，以使得泛型函数可以实例化。
	var i int = convertTo[int](i8)
	//var j int = convertTo(i8)
	_ = i
	f := func(s string) int {
		var result int
		result, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		return result
	}
	var strs = []string{"1", "2", "3", "4", "5"}
	var ints = Map(strs, f)
	var ints2 = []int{1, 2, 3}
	var strs2 = Map(ints2, strconv.Itoa)
	fmt.Println(ints)
	fmt.Println(strs2)

}

// //////////////////基于约束的类型推断/////////////////////
type integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

func Double[E integer](es []E) []E {
	r := make([]E, len(es)) //   注意，这里创建的是[]E
	for i, e := range es {
		r[i] = e + e
	}
	return r
}

// 这泛型函数中，类型参数S的约束或元类型（Meat—Type）是由另外的一个类型参数E所定义的，
// 所以，称S的约束或元类型是结构化约束（structural constrain）。
// 这个例子展示了如何使用类型参数定义其他的类型参数的约束，以及取得的效果。

func DoubleDefined[S ~[]E, E integer](s S) S {
	r := make(S, len(s)) // 注意，这里创建的是 S
	for i, v := range s {
		r[i] = v + v
	}
	return r
}
func TestTypeInferenceByConstrain(t *testing.T) {
	type MySlice []int
	//V1 的类型，也就是Doubl函数中的E类型被推断为int，因为常量1，2，3的缺省类型是int。
	//这里只用到了函数参数类型推断，没有用到基于约束的类型推断。
	//使用输入给函数的参数来推断类型，[]E -> MySlice=[]int,故而E的类型是int。
	var V1 = Double(MySlice{1, 2, 3}) //MySlice{1, 2, 3}被隐式转换为[]int{1, 2, 3}，传递给Double函数。
	fmt.Println(V1)
	//V2的类型，也就是DoubleDefined函数的返回值类型为MySlice。
	//这里函数参数的类型推断不能推断出所有的类型，因为类型参数E没有用在函数的输入参数中。
	//所以这里还要使用的是基于类型约束进行类型推断。
	// 首先，通过函数参数类型推断得出，S -> MySlice,
	// 又根据
	//S = ~[]E  -> MySlice =  []int ,故而E的类型是int，
	var V2 = DoubleDefined(MySlice{1, 2, 3}) //MySlice{1, 2, 3}被直接传递给DoubleDefined函数。
	fmt.Println(V2)

}

// //////////////////////////////////////////////////////////////////
// //////////////指针方法案例 Pointer Method Example////////////////////////////////////////////
// 这是一个类型约束
type Setter interface {
	Set(string)
}

// FromString函数使用字符串切片作为输入，返回类型T的切片。
// 注意，由于类型参数T没有用在函数的输入中，所以，无法使用函数输入参数进行类型推断。
// 如果 T是指针类型这个方法的实现会产生运行时错误
func FromStrings[T Setter](s []string) []T {
	result := make([]T, len(s)) //决定了T不能是指针类型。指针类型的空值是nil
	for i, v := range s {
		result[i].Set(v) // 如果T是指针，则reuslt[i]=nil,就会导致运行时的panic。
	}
	return result
}

type Settable int

// 注意是*Settable实现了 Set方法，而不是Settable实现了Set方法。
func (p *Settable) Set(v string) {
	if i, err := strconv.Atoi(v); err != nil {
		panic(err)
	} else {
		*p = Settable(i)
	}
}
func TestPointerExample1(t *testing.T) {
	var strs = []string{"1", "3", "5"}

	//var ss=FromStrings[Settable](strs) //Settable并没有实现Set方法。
	//编译oK，运行时，空指针panic
	//而且，我们希望得到的是 []Settable,不是[]*Settable
	var ss = FromStrings[*Settable](strs) //*Settable实现了了Set方法，但是*Settable的零值是nil，使用*nil导致运行时异常(panic)。
	fmt.Println(ss)
}

// 在Go语言中，普通接口是对“引用了具体类型”的绝对相同的方法签名的类型的一种抽象，可称之为“严格”共性行为的抽象。
// 而泛型接口，则是对“引用了可变类型参数”的语义相同的方法签名的类型的一种抽象，可称之为“宽松”共性行为抽象。
// 由于普通接口的“严格”共性要求，使得具体类型的签名必须与接口签名保持完全一致。
// 如果每个具体类型的共性行为都要操作自身类型，那这些具体类型就无法普通接口所要求的“严格”共性要求，
// 除非，每个具体类型都将共性行为中被操作的自身类型变为相同的接口类型，这就使得具体类型必须依赖抽象类型。
// 带来的问题就是抽象与具体分离的理念不够彻底，因为具体类型本身不应该知道抽象才对，因为它可能被从多个角度抽象。
// 泛型接口的出现使得抽象可以与具体彻底分离。泛型接口用“可变的类型参数”作为共性行为中被操作的类型，因此，放宽了
// 对共性行为的抽象范围，每当用一个具体类型来实例化一个泛型接口，就产生了一个新的，普通的，接口类型，
// 这个新的，普通的，接口类型则表达了“严格”的共性要求。而用来实例化该泛型接口的具体类型则往往是实现了
// 泛型接口共性行为的类型。
// Setter2是一个类型约束，它定义了这样一种类型：
// 1.实现了一个设置string值的Set方法。
// 2. 而且，是用这个类型的指针类型来实现Set方法。
// 总体来说，Setter2定义了一个类型集合，该集合中的类型是T的指针类型，并实现了Set(string)方法 。
type Setter2[T any] interface {
	Set(string)
	*T // 要求用T的指针类型满足此约束，也就是说用T的指针类型实现Set(string)方法。
}

//FromStrings2中引入了两个类型参数，T和PT，
//T用来定义返回切片的类型,即 []T,
//而Seeter2[T]约束所定义的参数类型表示一个指向了T的指针类型，*T，
// 也就是说，*T类型满足了Setter2约束的要求，实现了Set方法。

func FromStrings2[T any, PT Setter2[T]](strs []string) []T {
	result := make([]T, len(strs))
	for i, v := range strs {
		var pt PT = &result[i]
		pt.Set(v) // 如果T是指针，则reuslt[i]=nil,就会导致运行时的panic。
	}
	return result

}
func TestPointerExample2(t *testing.T) {
	var strs = []string{"1", "3", "5"}
	// 由于泛型的类型参数没有出现在函数的输入参数中，所以只能使用基于约束的类型推断，
	// 这里，PT Setter2[T] 表明类型参数PT的元类型（约束）是一个结构化约束。
	// 而，基于约束的类型推断从第一个类型参数开始推断，所以T -> Settable,
	// PT-> *T -> *Settable
	// *Settable实现了Set方法，且是指针类型，满足了Setter2[T]约束。
	var ss = FromStrings2[Settable](strs) //*Settable实现了了Set方法，但是*Settable的零值是nil，使用*nil导致运行时异常(panic)。
	fmt.Println(ss)
	var vStbl Settable = Settable(200)
	var pv *Settable = &vStbl
	_ = pv
	//含联合类型 (T1|T2|T3) 、底层类型(~T)、指针类型的（*T）等元素的接口不能
	// 当作普通接口使用，只能当作泛型参数的类型约束使用.
	//var st2 Setter2[Settable] = pv
	//_ = st2
}

// //////////////////////////////////////////////////////////////////////////////////////////////
// ////////////////类型集与方法共存的约束定义/////////////////////////////////////////////////////
type StringableSignedInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
	String() string
}
type MyInt64 int64
type MyInt32 int32

func (this MyInt32) String() string {
	return strconv.Itoa(int(this))
}
func PrintString[T StringableSignedInteger](t T) {
	println(t.String(), " is print by generic function ")
}
func TestComplextConstrain(t *testing.T) {

	var myInt32 = MyInt32(456)
	println(myInt32)
	//含联合类型 (T1|T2|T3) 、底层类型(~T)、指针类型的（*T）等元素的接口不能
	// 当作普通接口使用，只能当作泛型参数的类型约束使用.
	//var v StringableSignedInteger = myInt32
	PrintString(myInt32)
	var myInt64 = MyInt64(123)
	//MyInt64没有实现String() string方法，不符合StringableSignedInteger接口规范。
	println(myInt64)
	//PrintString(myInt64)

}

// ///////////////////////////////////////////////////////////////////////////////
// ////////////使用复合（composite）类型定义的约束///////////////////////////////////////////////////
// 对于使用复合类型所定义的约束，也就是约束所定义的类型集合中的类型都是复合类型，golang泛型对这样的约束（类型集合）
// 附加了一个要求：
// an operation may only be used if the operator accepts identical input types (if any)
// and produces identical result types for all of the types in the type set
// 对泛型的类型参数可使用的共性方法对于类型集合中的每个类型而言，都能够“操作同一种输入类型，而且返回同一种结果类型”。
type byteseq interface {
	//复合类型，也就是多个类型组成的类型，这里是数组类型[]与byte类型共同组成的复合类型.目前共8种。
	// 接口（interface），数组（array），指针（pointer）、结构（struct）、函数(function)、切片(slice)、映射（map）、通道（channel）都是复合类型。
	string | []byte
}

// Join函数将第数组a中的元素使用连接符 sep连接起来，并返回连接结果。
func Join[T byteseq](a []T, sep T) (ret T) {
	if len(a) == 0 {
		return ret
	}
	if len(a) == 1 {
		var emptyByts = []byte(nil)
		var firstBytes []byte = []byte(a[0]) //a[0]类型可能是[]byte或string，但是都可以转换为[]byte
		var resultBytes []byte = append(emptyByts, firstBytes...)
		return T(resultBytes)
		// 上面只是为了展示GO泛型能力在编译期间对类型参数的正确解析，相当于下面的代码。
		//return  a[0]
	}
	// 思路是头元素单独处理，然后把sep与剩余元素配对处理。
	n := len(sep) * (len(a) - 1) //内置len方法适用于包括string和[]byte在内的类型。
	for _, v := range a {
		n += len(v)
	}
	b := make([]byte, n)
	bp := copy(b, a[0]) //copy方法可以将string或[]byte拷贝到[]byte中。
	for _, v := range a[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], v)
	}
	return T(b)
	/**
	//上面代码主要讲述复合类型的泛型共性方法。下面代码才是简洁，其思路是：
	// 取出数组中的第一个元素，作为头元素，让sep与每个剩余元素配对，再组成一个元素，将这些元素拼在一起，
	// 就形成了最终结果。
	var resultBytes []byte
	resultBytes = append(resultBytes, a[0]...)
	for _, t := range a[1:] {
		var eleBytes = append([]byte(sep), t...)
		resultBytes = append(resultBytes, eleBytes...)
	}
	return T(resultBytes)
	**/
}
func TestCompositeTypeGeneric(t *testing.T) {
	var strs = []string{"hello", "World!", "I", "love", "you"}
	var sep = ","
	var s = Join(strs, sep)
	println(s)
}

// 注意，规范说可以将拥有相同的属性结构体进行抽象，定义为一种类型约束接口，
// 注意，并可以在泛型程序中操作满足此约束的具体结构体类型的属性。但到1.19尚不支持。
// 注意，已经支持对多个具有相同属性的结构体类型的约束接口类型的定义。
type SAX struct {
	A int
	X int
}

func (s *SAX) GetX() int {
	return s.X
}

type SBX struct {
	B int
	X int
}

func (s *SBX) GetX() int {
	return s.X
}

type SCX struct {
	C int
	X int
}

func (s *SCX) GetX() int {
	return s.X
}

type HasGetXMethod interface {
	GetX() int
}
type StructWithXfield interface {
	SAX | SBX | SCX
	GetX() int
}

//注意，在泛型类型中无法对具有相同属性的结构体类型参数的共性属性进行操作。
//IncrementX函数“不正确（INVALID）” .
// 实事上，就算是p.x的返回类型都一样，该函数还是不正确。
//可能是因为最终的Go语言规范(spec)与提案不一样，
//可能的原因是Go语言规范目前只考虑到类型集合只定义相同的方法集，而不是字段集。

func IncrementX[T StructWithXfield](p T) {

	v := p.GetX() //注意，目前还不支持
	v++
	p.X = v
}

// sliceOrMap is a type constraint for a slice or a map.

// 从提案来看，Entry函数中return c[i]应该允许，因为[]int与map[int]int两个类型都支持整数类型的下标操作，
// 并且都能返回同样的的类型
// 但是，实际上的Go语言规范则不允许这样做。因此，return c[i]不允许。
// https://stackoverflow.com/questions/71198899/how-to-constrain-type-to-types-with-index
// 提到了这一点：
// The type parameter proposal suggests that map[int]T could be used in a union with []T,
// however this has been disallowed.
// The specs now mention this in Index expressions:
// "If there is a map type in the type set of P, all types in that type set must be map types,
//
//	and the respective key types must be all identical".
//
// 翻译： 类型参数提案建议map[int]T 可以与[]T类型组成一个联合约束，但是这一建议被拒绝了。
// 规范在“Index expressions”中提到：
// 如果map类型出现在类型集合P中，那么类型集合P中的所有类型都必须是
type sliceOrMap interface {
	[]int | map[int]int
}

func Entry[T sliceOrMap](c T, i int) int {
	// This is either a slice index operation or a map key lookup.
	// Either way, the index and result types are type int.
	//return c[i]
	return 0
}

/////////////////////////////////////////////////////////////////////////////////////////
///////////////////类型集合中的类型参数（Type parameters in type sets）//////////////////////////////
/**
在约束接口的元素中，允许类型字面量引用约束的类型参数。
**/
// SliceConstraint 是一个类型约束，该约束中，类型字面量~[]T 引用了约束的类型参数T，
// 从而定义了一个与类型参数切片相匹配的类型集合。
type SliceConstraint[T any] interface {
	~[]T // 类型字面量，引用了类型参数T。
}

// Mapping函数使用某个类型的切片和一个变换函数作为输入参数，
//并且，返回的是将变换函数作用于输入切片的每个元素上所得到的结果切片。
//Mapping函数返回结果切片与它输入的切片参数的类型相同，即使类型是“被定义的类型（defined type）"也没有关系。

func Mapping[S SliceConstraint[E], E any](s S, f func(E) E) S {
	r := make(S, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

// MySlice是一个简单的，被定义的类型(defined type).
type MySlice []int

// DoubleMySlice takes a value of type MySlice and returns a new
// MySlice value with each element doubled in value.
func DoubleMySlice(s MySlice) MySlice {
	// The type arguments listed explicitly here could be inferred.
	v := Mapping(s, func(e int) int { return 2 * e })
	// Here v has type MySlice, not type []int.
	return v
}
func TestTypeParameterInConstrain(t *testing.T) {
	var ms MySlice = []int{1, 3, 5, 7}
	var ms2 = DoubleMySlice(ms)
	fmt.Printf("ms2 is : %v\n", ms2)

}

/////////////////////////////////////////////////////////////////////////////////////////
//////////////////类型转换（Type conversions）//////////////////////////////////////////
/**
如果类型集合中的所有类型之间都可以相互转换，那么在泛型中，就可以对类型参数使用显式的类型转换操作。
**/
type integers interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// 在泛型函数Convert中，将类型参数运用于显示的类型转换，
// 比如，To(from) heFrom（to），是因为To和From的类型集合都是整数类型，
// 而所有的整数类型之间允许类型转换（但会丢失精度数据）。
func Convert[To, From integers](from From) To {
	to := To(from)        //将类型参数To用于显式的类型转换
	if From(to) != from { //将类型参数From用于显式的类型转换
		panic("conversion out of range")
	}
	return to
}

func TestTypeConvert(t *testing.T) {
	var a int8 = 10
	var b int = Convert[int, int8](a)
	println(b)
	var c int = 100000
	var d int8 = Convert[int8](c)
	println(d)

}

////////////////////////////////////////////////////////////////////
//////////////// 无类型常量（Untyped constants）///////////////////////
/**
  如果无类型常量想要与类型参数进行操作，其前提是无类型常量可以转换为
  目标类型集合中的任意一个类型。
**/
type integers2 interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func Add10[T integers2](s []T) {
	for i, v := range s {
		s[i] = v + 10 // OK: 10 can convert to any integer type
	}
}

// This function is INVALID.
func Add1024[T integers2](s []T) {
	for i, v := range s {
		_ = i
		_ = v
		//	s[i] = v + 1024 // INVALID: 1024 not permitted by int8/uint8
	}
}

// //////////////////////////////////////////////////////////////////////////////////////
// ////////通过嵌入约束构成新的类型集合(Type sets of embedded constraints)////////////////////
// Addable 是支持 + 操作符的类型
type Addable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~complex64 | ~complex128 |
		~string
}

// Byteseq  是字节序列，支持索引操作： string 或 []byte.
type Byteseq interface {
	~string | ~[]byte
}

// AddableByteseq  是一个由嵌入约束组成的类型集合，该类型集合中的类型支持 + 操作符和索引操作。
// 这是一个同时能够满足Addable 和 Byteseq约束的一类类型。
// 换句话说，只有 ~string 类型才可以满足这个集合对类型的要求。
type AddableByteseq interface {
	Addable
	Byteseq
	//Addable|Byteseq // 多个类型联合（union）的语义是: 符合多个类型其中一个类型即可。
	//
	// 似乎二者很相像，但是对于使用约束的泛型函数而言差别很大，
	// 本质上，嵌入式接口（约束）的语义是构成方法的集合，也就是定义约束的方法集。
	//  而类型联合，在语义上则是直接构成类型集合。
	//二者的区别主要体现在：
	// 嵌入方式形成的约束对符合约束的类型集合要求苛刻（类型的交集），但对于作用于该约束所定义的操作和方法要求比较宽松，
	// 因为交集类型可以满足求交集的每类类型的所有操作方法（方法的并集）。
	// 联合方式形成的约束对形成的类型结合要求宽松（类型的并集），但是对作用于该约束所定义类型的操作或方法要求比较苛刻，
	// 只能使用所有类型的共同都存在的方法（方法的交集）。
}

func Add[T AddableByteseq](a, b T) T {
	return a + b
}

// 通过约束联合的方式形成新的约束或类型集合
type AddableOrByteseq interface {
	Addable | Byteseq
}

/**
//此函数无效，无法通过编译
func Plus[T AddableOrByteseq](a, b T) T {
	//return a + b //AddableOrByteseq要求对类型的操作必须是集合中所有类型的方法交集。[]byte类型不支持加操作。
}
**/
// 下面的函数有效
func PointTo[T AddableOrByteseq](a T) *T {
	return &a //AddableOrByteseq要求对类型的操作必须是集合中所有类型的方法和操作的交集。
}

func TestEmbbededConstrains(t *testing.T) {
	var a1 int = 10
	var b1 int = 12
	_ = a1 + b1 //int 支持+操作符
	//_=a1[1] //int不支持数组索引操作
	//Add(a1, b1) // int 类型不符合Byteseq约束。
	var a2 []byte = []byte{'a', 'b', 'c'}
	var b2 []byte = []byte{' ', 'g', 'o', 'o', 'd'}
	_ = a2[1] //[]byte类型支持数组索引操作
	_ = b2[1]
	//_ = a2 + b2 // []byte类型不支持+ 操作符
	//Add(a2, b2)// []byte类型不支持+操作符，不符合Addable约束

	var a3 string = "hello "
	var b3 string = "world"
	var char byte = a3[1] //string支持数组的索引操作
	_ = char
	_ = a3 + b3 // string支持 + 操作符。
	var s string = Add(a3, b3)
	println(s)
}

/////////////////////////////////////////////////////////////////////////////
///////////////类型联合中的接口类型（Interface types in union elements）////////////////////////////////////////
/**
  虽然Go的泛型提案中允许类型联合中可以使用任何类型，也包括普通接口类型，Stringish所示的那样，
  但是GO1.18语言规范尚不支持该提案。
**/
/**
type Stringish interface {
	string | fmt.Stringer
}
**/
///////////////////////////////////////////////////////////////////////////////////
/////////////////存在空类型集合的可能性（Empty type sets）///////////////////////////
//下面的Unsatisfiable约束就定义了空的类型集合
// Unsatisfiable is an unsatisfiable constraint with an empty type set.
// No predeclared types have any methods.
// If this used ~int | ~float32 the type set would not be empty.
type Unsatisfiable interface {
	int | float32
	String() string
}

func ToPointer[T Unsatisfiable](t T) *T {
	return &t
}

type MyString string

func (ms MyString) String() string {
	return string(ms)
}

func TestEmptyTypeSets(t *testing.T) {
	var i int = 10
	_ = i
	var f float32 = 12.2
	_ = f
	var ms MyString = "hello"
	_ = ms
	//以下方法都不可用
	//pi := ToPointer(i)
	//pf := ToPointer(f)
	//pms := ToPointer(ms)
}

// ////////////////////////////////////////////////////////////////////////////////////////////
// ////////////////无法对于大约类型约束的类型变量进行类型判断,无法精准识别类型（Identifying the matched predeclared type）/////////////////////
type Float interface {
	~float32 | ~float64
}

func NewtonSqrt[T Float](v T) T {
	var iterations int
	switch (interface{})(v).(type) {
	case float32: //不支持写： case ~float32
		iterations = 4
	case float64: //不支持写： case ~float64
		iterations = 5
	default: //对于符合Float约束的，以float32或float64为底层类型的类型会引起下面的panic
		panic(fmt.Sprintf("unexpected type %T", v))
	}
	// Code omitted.
	_ = iterations
	return v
}
func TestTypeAssert(t *testing.T) {
	type MyFloat float32
	var G = NewtonSqrt(MyFloat(64))
	fmt.Println(G)
}

/////////////////////////////////////////////////////////////////////////
///////////类型参数没有办法支持类型转换///////////////////
//Copy1函数中试图使用类型参数进行类型转换，现在无法支持类型参数之间保证安全的类型转换。
/**
func Copy1[T1, T2 any](dst []T1, src []T2) int {
	for i, x := range src {
		if i > len(dst) {
			return i
		}
		dst[i] = T1(x) // INVALID，类型参数不支持类型转换。
	}
	return len(src)
}
**/

// 实现上述的拷贝功能，只能使用不安全的类型转换，由函数的使用者保证类型之间可以相互转换。
func Copy2[T1, T2 any](dst []T1, src []T2) int {
	for i, x := range src {
		if i > len(dst) {
			return i
		}
		dst[i] = (interface{})(x).(T1) //  不安全的类型转换方式。
	}
	return len(src)
}

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////现在还不支持“类型的方法”中定义额外的类型参数/////////////////
/**
在GO语言中，类型的方法主要作用是让类型得以实现接口。由于抽象接口与具体类型之间可以比较彻底地分离，
因此，在类型方法中使用类型参数，该类型方法存在多种具体类型下的方法实例。
对于与接口分离比较彻底的Go语言来说，增加识别该具体类型是否是符合接口的困难。
目前，当前，可行的方案还没有明确。
下面是来自
https://go.googlesource.com/proposal/+/HEAD/design/43651-type-parameters.md#no-parameterized-methods
的例子，详细情况见原文。
**/
/**
// S is a type with a parameterized method Identity.
type S struct{}

// Identity is a simple identity method that works for any type.
func (S) Identity[T any](v T) T { return v } //类型S的方法。

type HasIdentity interface {
	Identity[T any](T) T  //接口中的函数都属于“类型的方法"。
}
**/

///////////////////////////////////////////////////////////////////////////////////////
//////////////// 没有办法使用参数类型对应的指针类型的方法//////////////////////////////////////////////
// Stringify2 calls the String method on each element of s,
// and returns the results.

type Stringer interface {
	String() string
}

func Stringify2[T Stringer](s []T) (ret []string) {
	for i := range s {
		ret = append(ret, s[i].String())
	}
	return ret
}
func TestUseTypeParmPointerMethods(t *testing.T) {
	var bf *bytes.Buffer = bytes.NewBuffer([]byte{'1', '2', '3', '4', '5'})
	var str string = bf.String()
	fmt.Println("result is :", str)
	//拷贝*bf的内容，形成数组
	var bfs []bytes.Buffer = []bytes.Buffer{*bf, *bf}
	_ = bfs
	//  Stringify2[bytes.Buffer](bfs) //无法通过编译，因为 String()方法属于 *bytes.Buffer类型
}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////
