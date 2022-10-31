package containertypes

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"unsafe"
)

// GO语言中，所有通过变量赋值形式进行的值拷贝，都是拷贝值的“直接值部（Direct Value Part）”
// 而SizeOf函数也是用于度量“直接值部（Direct Value Part）”的字节大小。故而，凡是内容存储在"底层值部"中的
// 容器类型的值，其SizeOf函数测量的内存占用都是一样的，只是“直接值部”所占用的内存大小，不能反应真实内存的占用。
// 数组的内存布局只有一个连续的内存段作为“直接值部（Direct Value Part）”，没有“底层值部（Underlying Value Part）”，
// 所以，数组中的数据内容都存储在“直接值部”内，数组的拷贝会拷贝源数组的内容。而SizeOf所测式的数组所占字节的大小
// 也是数组中 数据元素的字节数与数组长度的乘积。
// 而Slice、Map、Channel、function、Interface、String不仅有“直接值部”，也有底层值块，因此拷贝仅拷贝的是
// 直接值部，而对于这些类型的不同值使用SizeOf进行内存你大小测定，都是相同的值，也就是直接值部的固定字节数。
func TestArrayCopy(t *testing.T) {
	a1 := [4]int{1, 2, 3, 4}
	a2 := a1
	for i, v := range a2 {
		if v != a1[i] {
			t.Error("数组copy出错,源数组与目标数组值不相等")
		}
	}
	a1[0] = 100
	if a2[0] != 1 {
		t.Error("改变源数组的内容，影响了目标数组的内容")
	}

}
func TestStringCopy(t *testing.T) {
	//字符串会开辟两段内存，直接值部存储字符串的长度(len)，以及字符存储段的指针引用（*byte），
	//底层值部存储了字符串中的字符.所以，go中，（包括strings包中）对字符串变量的任何变更操作
	//都会导致新的底层值部的产生。所以本质上，位于字符串底层的字符序列是不可变的，不允许通过[index]方式来改变
	//指定位置的字符。其他任何变动字符串底层字符序列的操作都会导致新的内存开辟与复制，
	//因此，字符串的变更的效率不高。故而大量字符串的变更操作应尽量使用bytes.Buffer。或者
	//但是，字符串的拷贝只是拷贝“直接值部”，其实是一种“假拷贝”，因此，字符串的”拷贝“效率很高，不过
	//由于golang设计中，位于字符串底层的字符序列的不可变性，使得”假拷贝“不会引起数据的不一致性。

	s1 := "hello"
	ch := s1[0]
	//s1[0]='a' // 不允许通过index修改字符串中的字符。
	fmt.Printf("s1[0]是 %c \n", ch)
	directValuePartAddressOfs1_1 := &s1
	fmt.Printf("s1的直接值部地址：%v\n", directValuePartAddressOfs1_1)
	s2 := s1 //s1和s2的直接值部的内存地址不同，但内容相同，但实际上共同引用了同一个存储字符串内容的“底层值部”。
	fmt.Printf("s2的直接值部地址：%v\n", &s2)
	if s1 != s2 {
		t.Error("字符串s1与s2不同")
	}
	s3 := "hello"
	fmt.Printf("s3的直接值部地址：%v\n", &s3)
	if s1 != s3 {
		t.Error("字符串s1与s3内容不同")
	}

	s1 += " world" //  s1的底层值部变化为另外一段内存，内容是"hello World"，直接值部的内存地址不变，但内容被更新为新值。

	directValuePartAddressOfs1_2 := &s1
	fmt.Printf("变更内容后，s1的直接值部地址：%v\n", directValuePartAddressOfs1_2)
	if directValuePartAddressOfs1_1 != directValuePartAddressOfs1_2 {
		t.Error("字符串s1的直接支部地址发生了变化")
	}
}

// 字符串的比较不仅比较直接值部，还会比较底层值部。而且支持三种比较方式：
// 1. 用操作符比较，== ,!=,>,<,>=,<=
// 2. strings.compare(s1,s2)函数, 返回0表示相等，返回1表示s1>s2,0表示s1==s2,-1表示s1<s2
// 3. strings.EqualFold(s1,s2)函数，忽略大小写，比较字符串是否相等。
func TestStringCompare(t *testing.T) {
	s1 := "hello"
	s2 := "world"
	s3 := "WORLD"
	fmt.Printf("%v>%v: %v ,%v==%v: %v ,%v<%v: %v\n", s1, s2, s1 > s2, s1, s2, s1 == s2, s1, s2, s1 < s2)
	fmt.Printf("strings.compare(%v, %v)=%v \n", s1, s2, strings.Compare(s1, s2))
	fmt.Printf("strings.compare(%v, %v)=%v \n", s2, s3, strings.Compare(s2, s3))
	//打印结果：strings.compare(world, WorLd)=true
	fmt.Printf("strings.EqualFold(%v, %v)=%v \n", s2, s3, strings.EqualFold(s2, s3))

}
func TestSizeOf(t *testing.T) {
	//所有数据存储在“底层值部”的容器类型，虽然值的内容不同，但其unsafe.Sizeof()测得的内存占用都相同,
	//无法反应真是的内存占用。
	//字符串
	s1 := "123456"
	s2 := "12"
	fmt.Printf("字符串：unsafe.Sizeof(\"%s\")=%v,unsafe.Sizeof(\"%s\")=%v\n", s1, unsafe.Sizeof(s1), s2, unsafe.Sizeof(s2))
	fmt.Printf("字符串：len(\"%s\")=%v,len(\"%s\")=%v\n", s1, len(s1), s2, len(s2))
	//通过已有切片制作子切片，子片的长度是起止序号之差，容量是底层数组长度减去起始位置序号，也就是现有长度加上空余的容量。
	sl1 := []int{1, 2, 3, 4, 5}
	sl2 := []int{1, 2}
	fmt.Printf("切片：unsafe.Sizeof(%v)=%v,unsafe.Sizeof(%v)=%v\n", sl1, unsafe.Sizeof(sl1), sl2, unsafe.Sizeof(sl2))
	fmt.Printf("切片：len(%v)=%v,len(%v)=%v\n", sl1, len(sl1), sl2, len(sl2))
	//所有数据存储在“直接值部”的容器类型，比如数组，unsafe.Sizeof()可以测得真实内存占用。
	a1 := [5]int{1, 2, 3, 4, 5} // 5个整数占5*8=40 个字节内存
	a2 := [3]int{1}             //3个整数占3*8=24个字节的内存
	fmt.Printf("数组：unsafe.Sizeof(%v)=%v,unsafe.Sizeof(%v)=%v\n", a1, unsafe.Sizeof(a1), a2, unsafe.Sizeof(a2))
	fmt.Printf("数组：len(%v)=%v,len(%v)=%v\n", a1, len(a1), a2, len(a2))
}

// 切片与字符串类似，真实的数据都存储在“底层值部”中。切片的拷贝实际上是对相同的“底层值部”建立了
// 多个不同的引用（通过直接值部）。
// 和字符串一样的是，切片没有真正删除元素的操作，只有真正追加元素的操作，到追加元素超过了切片容量时，
//就会导致新的“底层值部” 的内存分配与拷贝，此时会得到一个新的切片来引用新的“底层值部”。
// 和字符串不同的是，切片允许通过下标操作改变指定位置的元素值，而通过下标改变切片元素的值时，
// 也不会引起“底层值部”重新进行内存分配，所以，对同一“底层值部”的不同引用都会看到被更新的元素值的变化。

func TestSliceOperation(t *testing.T) {
	//切片的“零值”的元素个数与最大可容纳元素个数均为0.
	var zeroSlice []int
	fmt.Printf("零值切片的元素个数为 %v ，最大可容纳的元素个数为 %v \n", len(zeroSlice), cap(zeroSlice))
	// 通过字面量方式创建的切片，其元素个数（长度）与最大容量相同，都是初始化的元素数量。
	//故而每次追加一个元素都会导致重新分配底层值部的内存。追加操作会导致长度增加了追加元素个数，而容量翻倍。
	s1 := []int{1, 2, 3, 4}
	fmt.Printf("切片：%v 的元素个数为 %v ，最大可容纳的元素个数为 %v \n", s1, len(s1), cap(s1))
	s2 := s1     // 浅拷贝，s2 与s1引用了相同的底层数组。
	s3 := s1[1:] //制作一个子切片，相当于得到了一个删除第一个元素的切片，但子切片s3与s1实际上引用了相同的底层数组。
	fmt.Printf("s3%v\n", s3)
	for i, v := range s2 {
		if v != s1[i] {
			t.Error("切片copy出错，源切片与目标切片引用的不是同一个底层数据")
		}
	}
	//不改变底层数组元素个数，只是更新某个元素的值，不会引起底层数组的内存分配与拷贝。

	s1[1] = 100
	//此时，s2,s3依旧与s1的引用相同的底层数组。
	fmt.Printf("s1[1]=%v,s2[1]=%v,s3[0]=%v \n", s1[1], s2[1], s3[0])
	if s2[1] != 100 || s3[0] != 100 {
		t.Error("通过元切片引用改变底层数据，引用同一底层数据的目标切面没有访问到被改变的数据")
	}
	//////追加1个元素，超过容量追加数据，长度加1（4+1=5）,而容量加倍（4*2=8）
	s1 = append(s1, 5) // Append函数在s1容量不够时会拷贝s1底层数据到新的内存段中，并在新的内存段追加元素，然后返回对这个新底层数据的切片引用。
	fmt.Printf("超过容量追加元素后的切片s1:%v的长度是%v，容量是%v\n", s1, len(s1), cap(s1))

	if len(s1) == len(s2) {
		t.Error("源切片底层数据个数发生了变化，拷贝切片也随之智能改变了！")
	}
	s4 := s1[1:]
	s5 := s1[2:4]
	fmt.Printf("子切片s4:%v的长度是%v，容量是%v\n", s4, len(s4), cap(s4))
	fmt.Printf("子切片s5:%v的长度是%v，容量是%v\n", s5, len(s5), cap(s5))
	s4[0] = 342
	fmt.Printf("S1[0]: %v\n", s1[0])
	fmt.Printf("S1[1]: %v\n", s1[1])
	fmt.Printf("S2[1]: %v\n", s2[1])
	fmt.Printf("S3[0]: %v\n", s4[0])
	///////////////////////////////////////////////////////////////
	//通过make函数创建的切片，可以指定元素个数（长度——len）与最大容量（cap）。
	//切片长度决定访问是否超界，切片容量决定是否重新分配底层值部（底层数组）的内存。
	///////////////////////////////////////////////////////////////
	sl1 := make([]int, 3, 5) //创建一个元素个数为3，3个元素都是int的零值0，容量为5的切片。
	sl1[0] = 1
	sl1[1] = 2
	sl1[2] = 3
	//sl1[3] = 4 //虽然没有超过容量，但是更新的元素超过切片当前的元素个数的长度，故而引起索引超界panic
	fmt.Printf("切片sl1:%v的长度是%v，容量是%v\n", sl1, len(sl1), cap(sl1))
	sl2 := sl1
	fmt.Printf("切片sl2:%v的长度是%v，容量是%v\n", sl2, len(sl2), cap(sl2))
	/////////////////////////////////////
	////注意，由于追加操作没有达到最大容量，不会导致底层值部重新分配内存，
	////这个位置值 (_intArray[3])会被下面的另外个引用（Sl2）在相同位置所追加的值覆盖。
	sl1 = append(sl1, 5)
	//sl2[3] = 4 //虽然引用同一个底层数组，但是不使用append操作，sl2的长度值不会加1，会引起索引超界panic
	sl2 = append(sl2, 4)
	fmt.Printf("切片sl1:%v的长度是%v，容量是%v\n", sl1, len(sl1), cap(sl1))
	fmt.Printf("切片sl2:%v的长度是%v，容量是%v\n", sl2, len(sl2), cap(sl2))
	sl1 = append(sl1, 6) //sl1的操作没有超过所引用的底层数组的容量，sl1与sl2仍引用相同的底层数组。
	sl1[0] = 100         //sl2[0]也应是100
	if sl2[0] != 100 {
		t.Errorf("期望sl1[0]的值100，而得到的值是%v\n", sl2[0])
	}
	//通过sl1应用对底层数组追加1个元素，使切片超过容量，引起底层数组重新内存分配，
	//此时，sl1引用新底层数组，sl2仍引用原有的底层数组。
	//此时，sl1的长度加1（5+1=6），容量翻倍（5*2=10）
	sl1 = append(sl1, 7) //sl1的操作超过容量，引起底层数组重新内存分配，sl1引用新底层数组，sl2仍引用原有的底层数组
	fmt.Printf("超过容量追加元素后的切片sl1:%v的长度是%v，容量是%v\n", sl1, len(sl1), cap(sl1))
	sl1[0] = 200       //sl1所引用的新底层数组的第一个元素被修改
	if sl2[0] != 100 { //sl2所引用的老的底层数组的第一个元素仍然维持原来的值
		t.Errorf("期望sl1[0]的值100，而得到的值是%v\n", sl2[0])
	}
	fmt.Printf("切片sl1[0]=%v ,sl2[0]=%v \n", sl1[0], sl2[0])
	fmt.Printf("切片sl2:%v的长度是%v，容量是%v\n", sl2, len(sl2), cap(sl2))
}

// 测试对切片追加的操作
func TestSliceAppend(t *testing.T) {
	s1 := make([]int, 3, 5)

	fmt.Printf("初始切片s1=%v  \n", s1)
	s1_2 := append(s1, 3) //当追加元素后，底层值部容量未超界时，对s1直接值部没有影响，没开辟新的底层值部，原底层值部同时被s1_2的直接值部所引用
	fmt.Printf("未容量追加元素后s1=%v  \n", s1)
	fmt.Printf("未容量追加元素后s1_2=%v  \n", 1_2)
	fmt.Printf("切片s1:%v的长度是%v，容量是%v\n", s1, len(s1), cap(s1))
	fmt.Printf("切片s1_2:%v的长度是%v，容量是%v\n", s1_2, len(s1_2), cap(s1_2))
	//
	s2 := append(s1, 3, 4, 5, 6) //当容量超界时，对s1直接值部没有影响，开辟了新的底层值部，被s2直接值部所引用
	fmt.Printf("超容量追加元素后，切片s1=%v  \n", s1)
	fmt.Printf("超容量追加元素后，切片s2=%v  \n", s2)
	fmt.Printf("切片s1:%v的长度是%v，容量是%v\n", s1, len(s1), cap(s1))
	fmt.Printf("切片s2:%v的长度是%v，容量是%v\n", s2, len(s2), cap(s2))
}

// 容器类型的字面量需要书写多个元素的字面量，如果每个元素的字面量都要书写元素的类型，就会非常麻烦，
// GO 编译器很贴心，只要容器类型定义中声明了所含元素的类型，在书写容器字面量时，元素字面量的类型就可以
// 省略。
func TestContainerValueLiteralSimplify(t *testing.T) {
	type language struct {
		name string
		year int
	}
	//将数组字面量赋值给变量X
	var X = [...]language{
		language{"C", 1972},
		language{"Python", 1991},
		language{"Go", 2009},
	}
	//数组变量的赋值可以简写为一下形式，
	//注意，容器字面量的类型中已经指定元素变量类型，就可以在元素值的字面量的书写上省略元素的类型。
	var X2 = [...]language{
		{"C", 1972},
		{"Python", 1991},
		{"Go", 2009},
	}
	fmt.Printf("%v\n%v\n", X, X2)
	//变量heads被赋值了一个切片字面量。注意，切片中，元素的两类型是4个byte元素的数组指针。
	var heads = []*[4]byte{
		&[4]byte{'P', 'N', 'G', ' '},
		&[4]byte{'G', 'I', 'F', ' '},
		&[4]byte{'J', 'P', 'E', 'G'},
	}
	//前面的例子表明，只要在容器字面量的类型中指定元素变量类型，就可以在元素值的字面量中省略元素的类型。
	//注意，这里甚至可以为指针两类的元素省略取地址操作。
	var heads2 = []*[4]byte{
		{'P', 'N', 'G', ' '},
		{'G', 'I', 'F', ' '},
		{'J', 'P', 'E', 'G'},
	}
	fmt.Printf("%v\n%v\n", heads, heads2)
	pInt := new(int)
	*pInt = 100
	pString := new(string)
	*pString = "hello"
	fmt.Printf("整型指针变量的内存大小为%v个字节,值为%v,字符串变量的指针的内存大小为%v个字节, 值为%v\n", unsafe.Sizeof(pInt), *pInt, unsafe.Sizeof(pString), *pString)
}

// 容器类型本质上都有序号（index or key），数组和切片的序号都是非负整数，且从0开始，到size-1 结束。
// 而map[K,V]的序号就是K。因此，在容器值的字面量中，可以通过序号可以指定特定位置的元素值，
// 其他元素的值则为默认的“零值”。
// 在容器字面量中通过序号指定特定位置的元素时，可以不按照固定顺序设定元素值，但是序号必须常量，不能是变量。
// 因为完整的容器字面量是由编译器在编译期间自动推断，编译期间只能识别常量或者常量表达式，或者unsafe包中
// 的函数（返回值）。
func TestContainerValueLiteralElementOrderSetting(t *testing.T) {
	/////下三个数组是等价的////////
	a0 := [6]int{1, 0, 4, 0, 16, 0}
	a1 := [6]int{0: 1, 2: 4, 4: 16}
	a2 := [6]int{2: 4, 4: 16, 0: 1}
	//数组值只有一个直接值部，所以一般情况下数组都能比较，即，比较每一个元素是否相等，
	//除非数组中的元素不可比较。
	fmt.Printf("a1==a0 ? %v , a1==a2 ? %v\n", a0 == a1, a1 == a2)
	/////////以下三个切片是等价的////////////////
	//切片值之间不支持比较操作（==，!=），无法通过编译，但切片值可以与nil进行比较。
	s0 := []int{1: 10, 3: 30, 5: 50}
	s1 := []int{3: 30, 5: 50, 1: 10}
	s2 := []int{0, 10, 0, 30, 0, 50}
	fmt.Printf("s0: %v\ns1: %v\ns2: %v\n", s0, s1, s2)
	///////以下三个map是等价的//////////////////
	m0 := map[string]int{"A": 1990, "C": 1960, "D": 2000}
	m1 := map[string]int{"C": 1960, "D": 2000, "A": 1990}
	m2 := map[string]int{"D": 2000, "A": 1990, "C": 1960}
	fmt.Printf("m0: %v\nm1: %v\nm2: %v\n", m0, m1, m2)
	///数组或切片的字面量中，表示元素序号必须是常量，即，命名常量，或者uint的字面量///
	var i int = 1
	_ = i
	const b int = 1
	//	sl := []int{i: 100} //无法编译通过，因为必须通过常量
	sl := []int{b: 100} //ok,b是命名常量
	_ = sl
}

// 切片与map等有“底层值部”的容器类型的“零值”是nil，表示该容器尚未创建，
//与已创建的但有0个元素的空容器完全不同，后者可以进行数据读写操作。
// 虽然切片值之间、容器值之间不可以比较，但是二者均可以与“空值（nil）”进行比较。
// 对零值切片(nil)求长度（len）与容量（cap），不会抛出运行时panic，得到的结果都是0,
// 但是，对零值切片(nil)进行元素数据的读写访问，都会抛出运行时panic，元素访问超界。
// 而对零值map(nil)进行数据元素的读取访问，不会抛出运行时panic，会得到对应类型的零值，
// 但是，对零值map(nil)进行数据元素的“写操作”，则会抛出运行时panic.
// 由于数组没有底层值部，所以数组“零值”不是nil，而是实实在在都是0值内存的直接值部，所以数组值之间可以比较，但不能与nil进行比较。

func TestContainerNilAndZeroValue(t *testing.T) {
	var a [16]byte                       //零值数组，分配了16个字节的内存，每个字节都是0.
	var s []int                          // nil，
	var m map[string]int                 //nil
	fmt.Println(a == a)                  // true
	fmt.Println(m == nil)                // true
	fmt.Println(s == nil)                // true
	fmt.Println(nil == map[string]int{}) // false，零值容器（nil）不等于“0个元素的空容器”
	fmt.Println(nil == []int{})          // false, 零值容器（nil） 不等于“0个元素的空容器”
}
func TestSliceValueAccess(t *testing.T) {

	s := []int{}                            //s!=nil
	fmt.Printf("%v,%v\n", len(s), cap(s))   // 0,0不会引起panic
	var s2 []int                            //s2==nil
	fmt.Printf("%v,%v\n", len(s2), cap(s2)) // 0,0不会引起panic
	i := s[0]                               //会产生运行时的panic，因为s=nil，s[0]会引起超范围访问数据
	fmt.Printf("%v\n", i)
}

// 对于map[K,V]，K必须是值可比较的类型。
// map的拷贝同slice拷贝一样，拷贝了直接值部，两个map变量通过直接值部共享一个底层值部，这样，任何一个
//map变量对元素的操作都会被另一个map变量所“看到“。与Slice不同的是，map的内存机制与slice底层值部使用
//连续的内存段不同，任何一个map变量对元素的操作都不会导致底层值部的重新内存分配，两个map变量会永远共享
//同一个底层值部。而slice的元素操作在超出容量后会导致底层值部的内存被重新分配，从而导致两个slice变量不再共享同一个底层值部。
// 使用K的值k对map进行取值时，当m为nil，或者m[k]的值不存在时，都不会抛出panic，
// 而会返回对应类型的零值（注意，有底层值部和无底层值部的零值不同，有底层值部的类型的零值是nil）.
// 同时,m[k]会返回两个值，第二个值的类型为bool,表示元素是否存在,第二个返回值可以忽略，故而
// m[k]即可以给一个变量赋值 ，也可以给两个变量赋值。

func TestMapValueAccess(t *testing.T) {

	m2 := map[string]int{"A": 123, "B": 234, "C": 456}
	i2, present := m2["A"]
	fmt.Printf("%v,%v\n", i2, present) //123,true
	i3, present := m2["NOVALUE"]
	fmt.Printf("%v,%v\n", i3, present) //0, falses
	m2["D"] = 678
	fmt.Printf("%v\n", m2["D"]) //678,
	//删除Map中的元素。由于slice底层值部是固定大小的内存块，与数组一样，
	//这种结构的元素内存不能被实际删除，只能刷新为零值，但可以通过子切片擦操作，
	//以”屏蔽视图“的方式实现删除的效果。
	//内置方法delete删除map元素不会产生任何错误，也不返回任何值，即使元素不存在，
	//甚至map是nil也不会有问题，只是做了一个空操作（ no-op）

	delete(m2, "NOVALUE") // key为"NOVALUE"的元素不存在，什么都不做
	m3 := m2
	m3 = nil
	delete(m3, "D")         //m3是nil，什么都不做
	println("ok,no panic!") // 可以打印
	//注意，这是容易产生bug的地方，
	//对值为nil的map进行取值，不会出现panic，而是返回值类型的空值。
	var m map[string]int = nil
	i := m["A"]
	fmt.Printf("%v\n", i) //0
	m["A"] = 222          //会抛出运行时异常，不能对nil map进行元素赋值操作。
}

// 内置make函数可以用于创建slice，map容器，但不能创建数组，同时还能创建特殊容器channel。
// make与new不同的是，make返回的是所创建容器的直接值部数据，可以赋值给容器类型的变量，
// new 返回的是所创建容器的直接值部的地址数据，可以赋值给容器指针类型的变量。
// 同时，new 不能指定容器的容量参数
func TestCreateContainer(t *testing.T) {
	s := make([]int, 3, 6)
	_ = s
	///使用内置函数make创建map
	m := make(map[string]int, 3)
	m2 := m
	m["A"] = 100
	m["B"] = 120
	m["C"] = 170
	m["D"] = 190
	for i := 0; i <= 100000000; i++ {
		m[strconv.FormatInt(int64(i), 10)] = i
	}
	fmt.Printf("%v\n", m2["D"])
	fmt.Printf("%v\n", m2["10000000"])
	//fmt.Printf("%v\n", m)
	////使用内置函数new创建map
	pm := new(map[string]int)
	fmt.Printf("%T\n", pm)

}

// 容器元素的可寻址性，是指是否可取得容器中元素的内存地址。
// 对于数组，如果数组是可寻址的，那么数组元素就可以寻址。数组字面量是不可以寻址的。
// 切片因为有底层值部，元素都在底层值部中，且切片元素不可删除，只可追加，所以可以对已存在元素寻址。
// 因此，任何切片中的元素都可以寻址，也就是说，尽管切片字面量自身不可寻址，但其中元素可寻置。
// map元素不可寻址，因为map元素可以删除，且可能重新分配内存，所以key所指定的元素位置不固定。

func TestContainerElementAddressability(t *testing.T) {
	var a = [...]int{1, 3, 5}
	pa0 := &a[0]
	*pa0 = 100
	s := make([]bool, 5, 10)
	ps2 := &s[2]
	*ps2 = true
	fmt.Printf("%v,%v\n", *pa0, *ps2)
	ps0 := &[]string{"Go", "C"}[0] //切片字面量的元素可寻址
	//pa0 := &[2]string{"Go", "C"}[0] //数组字面量的元素不可寻址
	fmt.Printf("%v\n", ps0)
	m := make(map[string]int, 5)
	m["A"] = 5
	//pma := &m["A"] //map元素不可寻址
	//fmt.Printf("%v\n", *pma)

}

// map条目中的值元素只允许整体更改，这对于大多数简单类型，以及含有底层值部的类型，比如切片等
// 都不会有影响，只是对不含底层值部，只有直接支付的复杂类型，比如数组和结构体来说会带来一定的不便。
func TestMapElementModify(t *testing.T) {
	type person struct {
		name string
		age  int
	}
	mp := map[string]person{}
	mp["kite"] = person{name: "kite", age: 20}
	//mp["kite"].age = 22 //无法编译，不允许更改map中结构体类型元素的局部数据
	ma := map[string][3]int{}
	ma["score"] = [3]int{100, 99, 98}
	//ma["score"][0] = 99 //无法编译，不允许更改map中数组类型元素的局部数据
	////////////切片类型元素可以更改局部元素 /////////
	ms := map[string][]int{}
	ms["score"] = []int{100, 99, 98}
	ms["score"][0] = 99
	fmt.Printf("%v\n", ms)
	/////////为了达到修改map中数组与结构体元素的局部值，必须将其赋值拷贝到临时变量中，然后再整体更新回去

	p := mp["kite"] //赋值拷贝
	p.age = 22
	mp["kite"] = p //更新map元素
	fmt.Printf("%v\n", mp)
}

// 制作子切片的操作，有两种形式
// 形式1：subSlice:=baseSliceOrArray[low:high]
// 形式2：subSlice:=baseSliceOrArray[low:high:max]
// 形式1相当于subSlice:=baseSliceOrArray[low:high:cap(baseSliceOrArray)]
// low表示新切片取自基础数据中起始元素序号，high表示截止元素序号，注意，截止序号不包括在内。
// max表示在基础数据中，最大截止的元素序号。(不包括在内)
// low,high,max 的取值范围：0 <= low <= high<=max <= cap(baseContainer)
// 取值范围表明，子切片的起止元素可以超过基础切片的末元素序号，但不是不能超过底层数据的容量范围内的合法序号。
// len=high- lwo ,cap=max-low
// 当low=0或high=len(baseSliceOrArray)时，可以省略。
func TestSubSliceOp(t *testing.T) {
	a := [...]int{0, 1, 2, 3, 4, 5, 6}
	s0 := a[:] // <=> s0 := a[0 :7:7] 数组a作为切片的底层值部，s0与a共享元素。
	fmt.Printf("a[...]=%v len=%v cap=%v\n", a, len(a), cap(s0))
	fmt.Printf("s0=a[:]=%v len=%v cap=%v\n", s0, len(s0), cap(s0))
	s1 := s0[:] // s1:=s0
	fmt.Printf("s1=s0[:]=%v len=%v cap=%v\n", s1, len(s1), cap(s1))
	s2 := s1[1:3] // <=>s2:=a[1:3:cap(s1)]不包括序号为3的元素，cap(s1)=7,所以，cap(s2)=7-1=6
	fmt.Printf("s2=s1[1:3]=%v len=%v cap=%v\n", s2, len(s2), cap(s2))
	s3 := s1[3:] // <=> s3 := s1[3:7]
	fmt.Printf("s3=s1[3:]=%v len=%v cap=%v\n", s3, len(s3), cap(s3))
	s4 := s0[3:5] // <=> s4 := s0[3:5:7] ，len=2 cap=4
	fmt.Printf("s4=s0[3:5]=%v len=%v cap=%v\n", s4, len(s4), cap(s4))
	//s4的长度为2，容量为4，虽然s5的长度设定为2，但容量设定为2，小于s4的容量，这样，一旦对s5追加数据，就不在与s4共享底层数组。
	s5 := s4[:2:2] // <=> s5 := s0[3:5:5]
	fmt.Printf("s5=s4[:2:2]=%v len=%v cap=%v\n", s5, len(s5), cap(s5))
	//虽然s4自身的长度是2，但是容量是4，故而，子切片操作的low与high的范围可是0——4
	s52 := s4[1:4] //len(s4)=2 cap(s4)=4, len(s6)=3,cap(s6)=3
	fmt.Printf("s52=s4[1:4]=%v len=%v cap=%v\n", s52, len(s52), cap(s52))
	s6 := append(s4, 77) //s4的容量为4，现有两个元素，故而s6与s4共享底层数组,进而与s0共享数组
	fmt.Printf("s6=append(s4,77)=%v len=%v cap=%v\n", s6, len(s6), cap(s6))
	fmt.Printf("s0=%v len=%v cap=%v\n", s0, len(s0), cap(s0))
	s7 := append(s5, 88) //len(s5)=2,cap(s5)=2,所以s7与s5不再共享底层数组，重新分配底层数组，长度加1，容量加倍，为4
	fmt.Printf("s7=append(s5,88)=%v len=%v cap=%v\n", s7, len(s7), cap(s7))
	s8 := append(s7, 66) //由于s7的容纳量为4，长度为3，所以s8与s7功能共享底层数组
	fmt.Printf("s8=append(s7,66)=%v len=%v cap=%v\n", s8, len(s8), cap(s8))
	s3[1] = 99
	fmt.Printf("s3=%v len=%v cap=%v\n", s3, len(s3), cap(s3))
	fmt.Printf("s0=%v len=%v cap=%v\n", s0, len(s0), cap(s0))
	fmt.Printf("数组a=%v\n", a)

}

// 切片你可以从数组构建，也就是以数组为底层值部，二者共享元素。
// 从go1.17后，也可以将切片变量转换为数组指针。要求数组类型的长度不能超过切片的长度（len）。
// 切片变量转换为数组后，二者仍然共享相同的元素。
func TestSliceConvertToArray(t *testing.T) {
	a := [...]int{0, 1, 2, 3, 4, 5, 6, 7}
	s0 := a[:]
	s0[0] = 100 //通过切片更改底层数组的数据
	fmt.Printf("数组a=%v\n", a)
	fmt.Printf("切片s0=%v\n", s0)
	s := []int{100, 200, 300, 400, 500, 600}
	pa1 := (*[6]int)(s) //pa1数组指针的数据访问范围与s数据访问范围完全一致。
	fmt.Printf("*pa=%v,type=%T\n", *pa1, pa1)
	//注意，下面从切片s的第2个元素开始构建元素个数为3的数组指针
	pa2 := (*[3]int)(s[1:]) //pa2指向s的第2个元素。
	pa2[0] = 999
	fmt.Printf("*pa2=%v,type=%T\n", *pa2, pa2)
	fmt.Printf("s=%v\n", s)
	type PA4 *[4]int
	pa4 := (PA4)(s) //ok
	_ = pa4
	//注意，数组类型的长度超过了切片的长度，会引发运行时的panic
	pa3 := (*[7]int)(s)
	fmt.Printf("*pa=%v,type=%T\n", *pa3, pa3)

}

// 内置copy函数可以完成从源切片向目标切片的数据复制。
// 二者之间的长度可以不同，该函数返回值是所复制元素的数量，为二者长度的最小值。
// 二者皆可以为nil，也不会抛出panic，只不过返回值为0
// 源切片与目标切片因为各自拥有不同的底层数组，所以才存在使用copy函数进行元素复制的必要。
func TestBuildInCopyFunc(t *testing.T) {
	a := [...]int{0, 1, 2, 3, 4, 5}
	s1 := a[:]
	var s2 []int = make([]int, 3)
	n1 := copy(s2, s1)
	s2[0] = 100 //s1不受影响，二者不共享底层数组
	fmt.Printf("s1=%v\n", s1)
	fmt.Printf("s2=%v\n", s2)
	fmt.Printf("从s2copy %v个元素到s1\n", n1)
	// 源和目标都是nil也不会panic
	var s3 []int //s3=nil
	var s4 []int // s4=nil
	n2 := copy(s4, s3)
	fmt.Printf("s3=%v\n", s3)
	fmt.Printf("s4=%v\n", s4)
	fmt.Printf("从s3copy %v个元素到s4\n", n2)

}

// 测试字符串你与字节切片你之间的关系：
// 1.字符串与字节切片[]byte 可以进行相互转换，这种转换的本质是数据赋值，二者不共享底层数据。
// 2. 字符串可以作为copy或者append的内置函数的第二个参数，来向目标字节切片拷贝或追加数据
func TestRelationBetweenStrAndByteSlice(t *testing.T) {
	str := "hello "     //字符是只读的
	sb := ([]byte)(str) //把字符串转化为字节切片,二者不共享相同的底层数据。
	sb[0] = 'H'         //对str没有影响
	fmt.Printf("str=%v,sb=%v\n", str, string(sb))
	str2 := string(sb) // 把字节切片转换为字符串，二者不共享底层数据
	sb[1] = 'E'
	fmt.Printf("str2=%v,sb=%v\n", str2, string(sb))
	// 语法糖用法，把字符串追加到字符切片中，注意，字符串后面要加...
	helloworld := append(sb, "world!"...)
	fmt.Printf("hellworld=%v\n", string(helloworld)) //把字节切片转换为字符串打印。

	// 把字符串拷贝到字符切片中，注意，不加字符串变量后面不加...
	copy(helloworld[:5], "HELLO")
	fmt.Printf("after copy hellworld=%v\n", string(helloworld))
}

// 测试对容器的遍历
// 使用 for  k,v:= range container 语法遍历容器中的元素。
// 注意，如果 k,v是在for循环中声明的，那么它们只是循环体中的临时变量，每次循环都会重新声明。
// k,v  的值通过赋值的方式拷贝而来，对k，v的修改对容器元素没有任何影响。
func TestContainerIterate(t *testing.T) {
	a := [...]int{0, 1, 2, 3, 4, 5, 6}
	//遍历数组，所有数组元素值加倍
	for i, v := range a {

		a[i] = 2 * v
	}
	fmt.Printf("[...]int=%v\n", a)
	//注意，对数组进行序号的边界范围划分操作可以得到切片。
	s := a[:] //对数组进行边界操作，得到一个以数组为的底层数据的切片。
	//遍历切片，所有切片元素值加倍
	for i, v := range a {

		s[i] = 2 * v
	}
	fmt.Printf("[]int=%v\n", s)
	m := map[string]int{"A": 1, "B": 2, "C": 3, "D": 4}
	fmt.Printf("map[string]int=%v\n", m)
	//遍历map，所有map元素的值加倍
	for k, v := range m {
		m[k] = v * 2
	}
	//删除所有map元素
	fmt.Printf("map[string]int=%v\n", m)
	for k, _ := range m {
		delete(m, k)
	}
	fmt.Printf("map[string]int=%v\n", m)
}

// 通过反射机制修改切片的长度与容量
// 设定的长度不能大于数组的容量，因一旦大于切片的容量，就可能访问到未知内存。
// 同样的道理，设定的新容量也不能大于数组的容量，但也不能小于数组的长度。
// 因为小于数组的长度会导致返回合法数据失败。
func TestModifyLenAndCapByReflect(t *testing.T) {
	s := make([]int, 2, 6)
	fmt.Println(len(s), cap(s)) // 2 6
	//通过反射修改切片的长度
	reflect.ValueOf(&s).Elem().SetLen(3)
	fmt.Println(len(s), cap(s)) // 3 6
	//通过反射修改切片的容量
	reflect.ValueOf(&s).Elem().SetCap(5)
	fmt.Println(len(s), cap(s)) // 3 5
}
func TestSlicClone(t *testing.T) {

	s := make([]int, 5, 10)
	s[0] = 1
	s[4] = 5
	fmt.Printf("s=%v,len=%v,cap=%v\n", s, len(s), cap(s))
	//方式1：go1.18后的空切片追加方式，容量余1。 注意，如果s为nil，或者长度为0，那么结果就是nil
	sClone1 := append(s[0:0:0], s...)
	sClone1[0] = 100
	fmt.Printf("sClone1=%v,len=%v,cap=%v\n", sClone1, len(sClone1), cap(sClone1))
	fmt.Printf("s=%v,len=%v,cap=%v\n", s, len(s), cap(s))
	//方式2：传统方式空切片追加方式，容量余1。
	_ = sClone1
	sClone2 := append([]int(nil), s...)
	sClone2[0] = 199
	fmt.Printf("sClone2=%v,len=%v,cap=%v\n", sClone2, len(sClone2), cap(sClone2))
	fmt.Printf("s=%v,len=%v,cap=%v\n", s, len(s), cap(s))
	//方式3：精准拷贝。由于append对空切片追加元素时，新容量算法按照容量(0)加倍无效，所以，新容量算法为所追加元素个数加1。
	//所以上面两种append的克隆方式会导致克隆数组的容量(cap)比拷贝的数据元素数量多1，可能浪费一点内存.
	//精准克隆的思路都是预先设计好的目标切片的长度和容量。
	sClone3 := make([]int, len(s)) //使sClone3的长度与源切片相同，并使sClone3的容量与长度相同。
	copy(sClone3, s)
	sClone3[0] = 299
	fmt.Printf("sClone3=%v,len=%v,cap=%v\n", sClone3, len(sClone3), cap(sClone3))
	fmt.Printf("s=%v,len=%v,cap=%v\n", s, len(s), cap(s))
	//方式4：更慢一点的精准拷贝。
	sClone4 := append(make([]int, 0, len(s)), s...)
	sClone4[0] = 399
	fmt.Printf("sClone4=%v,len=%v,cap=%v\n", sClone4, len(sClone4), cap(sClone4))
	fmt.Printf("s=%v,len=%v,cap=%v\n", s, len(s), cap(s))
	//对于方式3,4，如果s是零值（nil）切片，则克隆的不是零值切片，而是一个空切片
	//方式5, 绝对精准的克隆，如果s是零值（nil）切片，那么结果也是零值切片。
	var sClone5 []int
	if s != nil {
		sClone5 = make([]int, len(s))
		copy(sClone5, s)
	}
	sClone5[0] = 499
	fmt.Printf("sClone5=%v,len=%v,cap=%v\n", sClone5, len(sClone5), cap(sClone5))
	fmt.Printf("s=%v,len=%v,cap=%v\n", s, len(s), cap(s))
	///append与make+copy方式的效率比较.
	///在go1.15前，append方式比make+copy效率高，但在go1.15后，经过优化，
	///一种常用的全量拷贝方式下，make+copy比append方式效率高。但，这种优化必须遵循严格的书写方式,
	//即：源切片必须是纯的独立变量，不能包含在任何其他数据中，比如容器或结构体。
	// 同时，make函数只能设定目标切片的长度不能设定容量参数。

	//优化起作用的case
	var sClone6 []int
	if s != nil {
		sClone5 = make([]int, len(s)) //优化OK，正确的使用方式
		copy(sClone6, s)
	}
	//优化无效的case1
	if s != nil {
		sClone5 = make([]int, len(s), len(s)) //注意 优化不起作用，make函数设定了容量参数。
		copy(sClone6, s)
	}
	//优化无效的case2
	bytes := []byte{'a', 'b'}
	var a = [1][]byte{bytes}
	y := make([]byte, len(a[0])) //注意 优化不起作用，拷贝源来自于其他容器
	copy(y, a[0])
	//优化无效的case3
	type T struct{ x []byte }
	var t1 = T{x: bytes}
	z := make([]byte, len(t1.x)) // 注意 优化不起作用，拷贝源来自于结构体。
	copy(z, t1.x)

}

func TestDeleteSegmentFormSlice(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	s2, _ := delSegmentWay1(s, 3, 6)
	showSlice("s2", s2)

	s3, _ := delSegmentWay2(s, 3, 6)
	showSlice("s3", s3)
	s = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	s4, _ := delSegmentWay3(s, 3, 5)
	showSlice("s4", s4)
	showSlice("s", s)
	s = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	s5, _ := delSegmentAndClearTailElements(s, 3, 5, delSegmentWay3[int])
	showSlice("s5", s5)
	showSlice("s", s)

}

// 删除元素，保留顺序
// 得到的结果与源共享底层数据元素，实际上修改了底层数据的元素，并产生了一个新的“视图”
func delSegmentWay1[T any](s []T, from, to int) ([]T, error) {
	if from > to {
		return nil, errors.New("form  大于 to")
	}
	//注意，因为from<=to<=len<=cap,s[:from]的容量为cap-0=cap，足够追加s[to:]元素，
	//所以，不会产生新的内存分配，二者共享同一个底层元素，只不过被删除数据的内存处被合法数据覆盖了
	//
	result := append(s[:from], s[to:]...)
	return result, nil
}

// 删除元素，保留顺序
func delSegmentWay2[T any](s []T, from, to int) ([]T, error) {
	if from > to {
		return nil, errors.New("form  大于 to")
	}
	//把to及其以后未被删除的元素拷贝到from及以后的元素中，拷贝的元素数量与from之前的数量之和即是剩余数量。
	//注意，没有使用make新切片，这时使用拷贝是对同一个底层内存进行数据元素的复制。
	result := s[:from+copy(s[from:], s[to:])]
	return result, nil
}

// 删除元素，可能不会不保留原有元素之间的顺序，这样拷贝的元素数量最少，效率最高
func delSegmentWay3[T any](s []T, from, to int) ([]T, error) {
	if from > to {
		return nil, errors.New("form  大于 to")
	}
	//如果剩余的后半段元素数量比删除数量的少，就拷贝剩余的后半段到删除位置处
	if delCount := to - from; len(s)-to < delCount {
		copy(s[from:to], s[to:])
	} else { //如果剩余的后半段元素比删除元素数量多，就截取与删除数量相等的最后一段数据拷贝到删除处
		copy(s[from:to], s[len(s)-delCount:])
	}
	//取得剩余结果
	result := s[:len(s)-(to-from)]
	return result, nil
}

// 注意 以上的slice的删除都是删除元素位置上拷贝了未被删除的元素，这样，底层数组中
// 总会保留一部分重复的数据，虽然通过“新视图”进行操作不会引起业务逻辑问题，但是，
// 如果尾部冗余的数据如果是指针或引用类型的数据，不设置为0值，那么在此视图存活期间，
// 被这些冗余元素所引用的数据就不会被垃圾回收，因此就会导致一些内存的浪费（泄露）。
func delSegmentAndClearTailElements[T any](s []T, from, to int, delFunc func([]T, int, int) ([]T, error)) ([]T, error) {
	result, err := delFunc(s, from, to)
	if err != nil {
		return nil, err
	}
	tobeClear := s[len(result):len(s)]
	var t0 T //T的零值
	for i, _ := range tobeClear {
		tobeClear[i] = t0
	}
	return result, err
}
func showSlice[T any](varName string, s []T) {
	fmt.Printf("%s=%v,len=%v,cap=%v\n", varName, s, len(s), cap(s))
}
func TestDeleteOneElementFromeSlices(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	showSlice("s", s)
	s1 := deleteElementWay1(s, 2)
	showSlice("s1", s1)
	showSlice("s", s)
	s2 := deleteElementWay2(s, 2)
	showSlice("s2", s2)
	showSlice("s", s)
	s3 := deleteElementWay3(s, 2)
	showSlice("s3", s3)

	s = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	showSlice("s", s)
	s4 := delElementAndClearTailElements(s, 1, deleteElementWay1[int])
	showSlice("s4", s4)
	showSlice("s", s)
	s = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	showSlice("s", s)
	s5 := delElementAndClearTailElements(s, 1, deleteElementWay3[int])
	showSlice("s5", s5)
	showSlice("s", s)
}

// 保持原有顺序方式1
func deleteElementWay1[T any](s []T, pos int) []T {
	if pos >= len(s) || pos < 0 {
		return s
	}
	result := append(s[:pos], s[pos+1:]...)
	return result
}

// 保持原有顺序方式2
func deleteElementWay2[T any](s []T, pos int) []T {
	if pos >= len(s) || pos < 0 {
		return s
	}
	result := s[:pos+copy(s[pos:], s[pos+1:])]
	return result
}

// 不保持原有顺序方式，最快的方式，把最后的元素填充到删除位置
func deleteElementWay3[T any](s []T, pos int) []T {
	sLen := len(s)
	if pos >= sLen || pos < 0 {
		return s
	}
	s[pos] = s[sLen-1]
	return s[:sLen-1]
}
func delElementAndClearTailElements[T any](s []T, pos int, delFunc func([]T, int) []T) []T {
	sLen := len(s)
	if pos >= sLen || pos < 0 {
		return s
	}
	result := delFunc(s, pos)
	var t0 T //T的零值
	s[sLen-1] = t0
	return result
}
func TestDeleteSliceElementsConditionally(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	var isKeep = func(v int) bool {
		return v > 5
	}
	showSlice("s", s)
	s1 := DelSliceElementsConditionally(s, isKeep, false)
	showSlice("s1", s1)
	showSlice("s", s)
	s = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	showSlice("s", s)
	s2 := DelSliceElementsConditionally(s, isKeep, true)
	showSlice("s2", s2)
	showSlice("s", s)
}
func DelSliceElementsConditionally[T any](s []T, keep func(T) bool, clear bool) []T {
	result := s[:0]
	for _, t := range s {
		if keep(t) {
			result = append(result, t)
		}
	}
	if !clear {
		return result
	}
	var t0 T                      //声明T的零值变量
	sTobeClear := s[len(result):] //获得底层数组中保留的无意义数据
	for i, _ := range sTobeClear {
		sTobeClear[i] = t0 //清0，防止内存泄露
	}
	return result
}

// 在指定位置插入另外一个切片中的元素，方式1，代码简单，但是可能效率低一点
func insertElementsWay1[T any](s []T, i int, elements []T) []T {
	//注意，两次append操作可能会导致两次内存的重新分配，所以效率会低。
	return append(s[:i], append(elements, s[i:]...)...)
}

// 在指定位置插入另外一个切片中的元素，方式1，代码啰嗦，但效率高一点
// 主要是通过容量判断来决定是否使用make方法手动构建新的，容量足够的切片。
func insertElementsWay2[T any](s []T, i int, elements []T) []T {
	sLen := len(s)
	eleLen := len(elements)
	sCap := cap(s)
	var rs []T
	// s容量足够
	if sCap >= sLen+eleLen {
		rs = s[:sLen+eleLen]        //注意，只要容量够，衍生切片的数据长度可大于主切片
		copy(rs[i+eleLen:], rs[i:]) // 把i及以后的元素拷贝（平移）到i+eleLen位置处
		copy(rs[i:], elements)      //把插入元素拷贝到i及后续位置

	} else { //s容量不足
		rs = make([]T, 0, sCap)
		rs = append(rs, s[0:i]...)   //把i之前的元素追加到新切片中
		rs = append(rs, elements...) //把插入元素追加到新切片中
		rs = append(rs, s[i:]...)    //把i及之后的元素追加到新切片。
	}
	return rs

}
func TestInsertAllElementsFormOtherSlice(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5, 6}
	ints := []int{100, 200, 300}
	s2 := insertElementsWay1(s, 2, ints)
	showSlice("s2", s2)
	showSlice("s", s)
	s3 := insertElementsWay2(s2, 1, ints)
	showSlice("s2", s2)
	showSlice("s3", s3)

}

func popFront[T any](s []T) (T, []T) {
	rt, rs := s[0], s[1:]
	//对被弹出的元素进行清零，否则容易有内存（临时）泄露
	var t0 T
	s[0] = t0
	return rt, rs
}
func popBack[T any](s []T) (T, []T) {
	rt, rs := s[len(s)-1], s[:len(s)-1]
	//对被弹出的元素进行清零，否则容易有内存（临时）泄露
	var t0 T
	s[len(s)-1] = t0
	return rt, rs
}
func pushFront[T any](s []T, t T) []T {
	return append([]T{t}, s...)
}
func pushBack[T any](s []T, t T) []T {
	return append(s, t)
}

// 只能popBack，pushBack的是堆栈。
// 只能popFront，pushBack的是队列
func TestSimulateStackWithSlice(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5, 6}
	showSlice("s", s)
	i, s := popFront(s)
	fmt.Printf("pop front :%v\n", i)
	showSlice("s", s)
	i, s = popBack(s)
	fmt.Printf("pop back :%v\n", i)
	showSlice("s", s)
	s = pushFront(s, 0)
	showSlice("s", s)
	s = pushBack(s, 6)
	showSlice("s", s)
}

// 使用map模拟Set
// struct{}类型的值不占用任何内存,可以用map[T]struct{}模拟Set
type Set[T comparable] map[T]struct{}

func (s *Set[T]) delete(t T) {
	delete(*s, t)
}
func (s *Set[T]) add(t T) {
	(*s)[t] = struct{}{}
}
func (s *Set[T]) isInclude(t T) bool {
	_, isInclude := (*s)[t]
	return isInclude
}
func (s *Set[T]) Iterate() []T {
	r := make([]T, 0, len(*s))
	for k, _ := range *s {
		r = append(r, k)
	}
	return r
}
func MapTo[T comparable, E comparable](fn func(T) E, st Set[T]) Set[E] {
	result := NewSet[E](0)
	for _, t := range st.Iterate() {
		e := fn(t)
		result.add(e)
	}
	return result
}
func NewSet[T comparable](startingSize int) Set[T] {
	temp := make(map[T]struct{}, startingSize)
	return (Set[T])(temp)
}
func TestSimulateSetWithMap(t *testing.T) {
	st := NewSet[int](10)
	st.add(1)
	st.add(2)
	st.add(3)
	st.add(4)
	for i, v := range st.Iterate() {
		fmt.Printf("st[%v]=%v\n", i, v)
	}
	st.delete(3)
	fmt.Println("----------after delete 3---------")
	for i, v := range st.Iterate() {
		fmt.Printf("st[%v]=%v\n", i, v)
	}
	fmt.Println(st.isInclude(3))
	fmt.Println(st.isInclude(4))
	helloInt := func(i int) string {
		return fmt.Sprintf("Hello %v", i)
	}
	st2 := MapTo(helloInt, st)
	for i, v := range st2.Iterate() {
		fmt.Printf("st2[%v]=%s\n", i, v)
	}

}
