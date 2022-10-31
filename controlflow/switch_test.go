package controlflow

import (
	"fmt"
	"testing"
)

/**
switch，就是开关的意思，也就是存在多个逻辑分支的情况下，通过与匹配目标进行对比，
决定到底执行“哪一个”或“哪几个”分支。
Go中的switch比C中switch用途更广泛，功能更强大，
主要实现思路是用switch 表达式作为匹配目标表达式，每个case分支都有“一个或者多个逗号分割”的条件表达式，
如果case分支的条件分支表达式中任何一个表达式的求值结果与switch匹配目标表达式求值结果相等，就意味着该分支匹配，
就会执行该分支下的代码。
具体语法是：
 switch  sw_exp {
	case c_exp1:
		case1 code
		...
		[fallthrough]
	case c_exp2:
		case2 code
		....
		[fallthrough]
	default:
		default case code
}

!!!程序运行时，先对switch 匹配目标表达式sw_exp求值，得到一个结果，也就是匹配目标值。
然后再按照顺序，对所有case分支的条件表达式分别进行求值，也就是对分支条件值。所有这些
分支条件都完成后，在自上而下，寻找分支条件值与匹配目标值相等（也就是匹配）的分支，然后执行该分支代码，
!!!如果在该分支最后没有fallthrough命令，就自动break，结束整个switch。
（这与C++的switch不同，c++必须写break才不会进入下一个分支代码）。
如果一直没有遇到匹配的case分支，但存在default分支，就执行该default分支。
如果相匹配分支的最后一条分支代码语句是fallthrough命令，那么不会自动结束switch，
而是忽略下个分支的表达式的求值与匹配，继续“直接、执行”下一个分支代码。
case分支代码中的fallthrough命令主要用于上个分支条件符合或满足下个分支条件，
而下个分支条件不满足上个分支条件的逻辑情况。这样，当上一个
分支条件可以与匹配目标相匹配时，下一个分支条件一定会满足匹配目标，因此，省略下个分支条件的求值与匹配将会提升执行效率。
比如，金牌会员肯定有银牌会员的特权 ，银牌会员肯定有普通会员的特权，普通会员肯定有游客的特权，反之则不成立。
此业务逻辑使用fallthrough，使得金牌会员全部拥有的特权构建逻辑从金牌特有分支一直降落至游客。
需要注意的是：

1.sw_exp必须是可以求出值（结果）的表达式，不能是变量、函数的声明、也不能是赋值语句等。
2.如果没有sw_exp表达式，则认为sw_exp表达式表达式求值结果为true，则 c_exp必须是一个布尔表达式，只要表达式求值为真，即表示匹配成功。
3.如果case表达式含函数调用，那么对表达式的求值就是对函数执行。
4.被求值case分支表达式个数，一定小于或等于case分支总数（当不存在匹配分支，或者最后一个分支匹配时二者数量相等）。
**/

//TestNormalSwitch演示了常规的Switch使用情况，也就是自顶向下匹配switch的表达式。

func TestNormalSwitch(t *testing.T) {
	peoples := [3]string{"jone", "kite", "lari"}
	findPeole := func(i int) string {
		fmt.Printf("get peoples[%v]: %v\n", i, peoples[i])
		return peoples[i]
	}

	var name string = "kite"

	switch name { //本例中，switch表达式为name变量,求值结果是kite，因此，各case分支要对比匹配的目标结果就是"jone"。
	case findPeole(0): //本例中，这个表达式会被求值，findPeole(0)函数被调用，但结果"jone"不匹配，导致下个case表达式被求值。
		fmt.Println("yes ,kite is the first people! ")
	case findPeole(1): //本例中，这个表达式会被求值，findPeole(1)函数被调用，分支求值结果与目标结果相等,都是"kite"，因此，该分支下的代码被执行，然后整个switch语句结束。
		fmt.Println("yes ,kite is the second people! ")
	case findPeole(2): //本例中，因为上个分支与目标结果匹配，整个switch语句结束，这个分支表达式不会被求值。
		fmt.Println("yes ,kite is the third people! ")
	default: //在本例中，第二个分支与目标结果相匹配，所以缺省分支不会被执行。
		fmt.Println("not find kite! ")
	}
}

// switch语句没有表达式，则switch求值结果为true,那么每个case语句的表达式求值结果必须bool，才能进行匹配
// 此时，相当于if...else...if...else

func TestNoSwitchExpression(t *testing.T) {

	testNoSwitchExpression := func(count int, weight int) {
		switch { //switch 语句后面没有表达式，那么分支匹配的目标结果就是true。
		case count < 100 || weight < 1000: //case 表达式求值结果必须是bool类型
			fmt.Printf("small packages !\n")
		case count > 10000 || weight > 100000: //case 表达式求值结果必须是bool类型
			fmt.Print("large packages !\n")
		default:
			fmt.Printf("midlle packages!\n")
		}
	}
	testNoSwitchExpression(10, 500)
}

// 在case中使用逗号(,)分割多个表达式，则会依次对每表达式求值，然后匹配，直到得到匹配项为止。
// 此种用法相当于把多个case子句合并，用一个case子句表达。
func TestMultiConditonInOneCase(t *testing.T) {
	peoples := [3]string{"jone", "kite", "lari"}
	findPeole := func(i int) string {
		fmt.Printf("evaluate peoples[%v]: %v in case statement\n", i, peoples[i])
		return peoples[i]
	}
	makePeole := func(name string) string {
		fmt.Printf("evaluate people: %v in switch statement\n", name)
		return name
	}
	var peopleName string = "kite"
	switch makePeole(peopleName) { //首先对switch表达式求值，也就是执行对makePeole函数的调用。
	case findPeole(0), findPeole(1), findPeole(2): //按照逗号逐个表达式求值匹配，只要有任何一个表达式匹配成功，则认为该分支匹配，就不再继续对后面的表达式求值。
		fmt.Printf("yes,find %v\n", peopleName)
	default: //如果上面的case都没有匹配成功，则执行
		fmt.Printf("not find the %v \n", peopleName)
	}
}

/*
*这个函数主要是为了理解如何正确使用GO语言中的类型断言，以便为通过switch进行类型判断做铺垫。

	类型断言函数返回两个值，第一个值类型转换结果，第二个值是bool类型，用于指明是否类型可以转换。
	变量类型断言的使用法在go语言中有点特殊，主要有以下三点：
	1.类型断言是系统函数，其返回结果不能全部抛弃。
	2.如果只抛弃第一个结果，也就是抛弃了类型转换结果，那么即使类型不匹配，无法转换，也不会抛出panic。
	3.如果只抛弃第二个结果，将第一个结果也就是转换结果赋值给某个变量，那么当类型不匹配的时候就会抛出panic。
	4.如果两个结果都不抛弃，那么即使类型不匹配，也不会抛出panic，此时，由于类型不匹配，第一个结果，
	 也就是类型转换的结果是目标类型的零值。

*
*/
func TestNoramlTypeAssertion(t *testing.T) {
	var unkonwTypeValue interface{} = "hello"

	//如果类型断言的两个值都不被抛弃，那么即使类型不匹配时，类型断言语句也不会抛出panic
	var i int = 8 //i初始化为非0值。
	i, ok := unkonwTypeValue.(int)
	if !ok {
		println("can not convert string  ‘hello’ to int，so i get ", i) //类型不匹配，i值为0（由于断言函数返回目标类型的0值）。
	} else {
		println(i)
	}
	////第一个参数被抛弃，即使类型不匹配，也不会抛出panic
	_, ok = unkonwTypeValue.(int)
	_ = ok
	//如果类型断言的最后一个值被抛弃，一旦类型不匹配，则会抛出panic
	i2 := unkonwTypeValue.(int) //类型不匹配，会抛出panic。
	_ = i2

}

// 将switch 语句用于接口变量的真实类型判断。
// type switch（类型判断分支）是switch语句的一种特殊用法。
// 只有在switch目标结果表达式语句中才能使用“类型断言表示式 var.(type)” 并赋值.
// 类型断言 var.(type) 不能用在switch目标之外的地方。
// 另外，只有类型断言表达式的赋值语句可以作为switch匹配目标表达式，其他赋值语句不能用作switch的匹配目标表达式。
// 接口变量的值是一个对具体类型值进行了"装箱(boxing)"的值，接口值通过存储具体类型值的指针，
// 以及具体类型来实现 “装箱(boxing)”,因此 switch 的obj.(type)表示中，匹配目标是type,
// obj.(type)的求值结果仍是接口类型的表达式，但是当type确定后，就可以"拆箱"为具体类型的值,
// 可以说var.(type)的求值结果本质上还是一个无法求出具体值的“拆箱表达式”。
func TestTypeAssertion(t *testing.T) {

	type Student struct {
		ID   string
		Name string
	}
	//为了避免面类型断言在运行时抛出异常，需要使用switch对类型进行判断，进行类型匹配的类型断言。
	typeAssertion := func(obj interface{}) {
		//objtype := obj.(type)    //obj.(type) 类型断言表达式只能用在switch中。
		//var i int;
		//switch i=1 { // 普通赋值语句也不能用在switch语句的匹配目标表达式与分支条件表达式中，因为赋值语句不返回值。
		//case 0:
		//	i=i+2
		//}
		switch t := obj.(type) { //obj.(type) 类型断言表达式只能用在switch中，变量t的类型为interface{}，本质是一个拆箱表达式。
		case int:
			var i int = t //此时，变量t所表达的拆箱表达式被求值为具体类型(int)的值。t相当于obj.(int)
			fmt.Printf("%v type is %T\n", i, i)
		case string: //此时，变量t所表达的拆箱表达式被求值为具体类型(string)的值。
			var s string = t //这里，t相当于obj.(string)
			fmt.Printf("%v type is %T\n", s, s)
		case *int: //此时，变量t所表达的拆箱表达式被求值为具体类型(*int)的值。
			var pi *int = t //这里，t相当于obj.(*int)
			fmt.Printf("%v type is %T\n", pi, pi)
		case *string: //此时，变量t所表达的拆箱表达式被求值为具体类型(*string)的值。
			var ps *string = t //这里，t相当于obj.(*strinng)
			fmt.Printf("%v type is %T\n", ps, ps)
		case Student: //此时，变量t所表达的拆箱表达式被求值为具体类型(Student)的值。
			var sd Student = t //这里，t相当于obj.(Student)
			fmt.Printf("%v type is %T\n", sd, sd)
		case *Student: //此时，变量t所表达的拆箱表达式被求值为具体类型(*Student)的值。
			var psd *Student = t //这里，t相当于obj.(*Student)
			fmt.Printf("%v type is %T\n", psd, psd)
		default:
			fmt.Printf("%v type is unkonwn\n", t)
		}
	}
	var i int = 10
	var str string = "hello"
	var std Student = Student{ID: "001", Name: "liufei"}
	var f float32 = 12.2
	typeAssertion(i)
	typeAssertion(&i)
	typeAssertion(str)
	typeAssertion(&str)
	typeAssertion(std)
	typeAssertion(&std)
	typeAssertion(f)
}

// fallthrough关键字用于switch case分支代码中，其作用恰恰与break中断switch语句相反，用于指定忽略
// 下个分支的表达式的求值与匹配，继续“直接、执行”下一个分支代码。
// fallthrough通常用于上一个分支条件兼容下一个分支条件，但下一个分支条件不兼容
// 上一个分支条件的情形。比如，一个整数大于100，自然大于10，但是大于10，未必大于100。
// 注意三点：
// 1. fallthrough会导致下一个分支不需要求值与匹配就会执行其分支代码。
// 2. fallthrough所在分支被执行才会产生效果。
// 3. fallthrough应运用在分支匹配条件兼容，合乎逻辑的情况下，否则容易造成误解与逻辑混乱。
// 本在示例中，设定每30年一个年龄段。
// 判断一个人的年龄为105的人，其年龄在90以上分枝，当然也满足60以上分枝、30以上分枝，以及0岁以上分枝。
func TestFallthrough(t *testing.T) {
	agePhases := [5]int{120, 90, 60, 30, 0}
	findAgePhase := func(i int) int {
		fmt.Printf("get age pahses[%v]: %v\n", i, agePhases[i])
		return agePhases[i]
	}

	const age int = 105
	println("the man's age is ", age)
	switch { //本例中，switch表达式为name变量,求值结果是kite，因此，各case分支要对比匹配的目标结果就是"jone"。
	case age > findAgePhase(0): //本例中，这个表达式会被求值，findPeole(0)函数被调用，但结果不匹配(105>200)，导致下个case表达式被求值。
		fmt.Println("his age is greater than 120  ") //不会被执行，
		fallthrough                                  //不会被执行。如果被执行，则会忽略下个分支的求值与匹配计算，而直接执行下个分支代码。
	case age > findAgePhase(1): //本例中，这个表达式会被求值，findPeole(1)函数被调用，结果匹配 (105>100),因此，该分支下的代码被执行。
		fmt.Println("his age is greater than 90 ") //被执行
		fallthrough                                //被执行。如果被执行，则会忽略下个分支的求值与匹配计算，而直接执行下个分支代码。
	case age > findAgePhase(2): //本例中，由于上个分支的代码中执行fallthrough，导致这个分支代码直接执行,而忽略了分支表达式的求值。
		fmt.Println("his age is greater than 60 ") //　上个分支代码中的fallthrough导致了本分支代码的直接执行。
		fallthrough                                //被执行。  如果被执行，则会忽略下个分支的求值与匹配计算，而直接执行下个分支代码。
	case age > findAgePhase(3): //本例中，由于上个分支的代码中执行fallthrough，导致这个分支代码直接执行,而忽略了分支表达式的求值。
		fmt.Println("his age is greater than 30 ") //　上个分支代码中的fallthrough导致了本分支代码的直接执行。
		fallthrough                                //被执行。  如果被执行，则会忽略下个分支的求值与匹配计算，而直接执行下个分支代码。
	case age > findAgePhase(4): //本例中，由于上个分支的代码中执行fallthrough，导致这个分支代码直接执行,而忽略了分支表达式的求值。
		fmt.Println("his age is greater than 0 ") //　上个分支代码中的fallthrough导致了本分支代码的直接执行。
	default:
		fmt.Println("his age is not right! ")
	}
}

// break在for和switch语句中的用途。
// break语句能够用来中断switch语句和for语句。
// 问题在于，当switch语句和for语句各自或相互多层嵌套时，break到底终止的是哪个呢？
// 存在for 、switch各自或相互多层嵌套的情况下，如果不加以其他处理，break遵循就近中断原则，
// 也就是中断break语句所在的for或switch。
// 如果想要中断某个外层的for或switch怎么办？此时需要对该外层的for或者switch语句定义一个标签,比如： label1，
// 然后使用break label1  来中断该层的for循环或switch。

func TestBreakStatementfunc(t *testing.T) {
loop: //为下面对变量i的for循环定义一个标签，便于break中断。
	for i := 0; i < 100; i++ {
		println("i loop to ", i)
	cases_select: //为下面的switch语句定义了一个标签，便于break中断。
		switch i {
		case 10:
			println("now, i is   10")
			for j := 0; j < 20; j++ {
				if j > 5 {
					break //此break中断的是对j的for循环，所在的for循环。
				} else {
					println("j  is ", j)
				}
			}
		case 30:
			println("now ,i is 30")
			for k := 0; k < 20; k++ {
				if k > 5 {
					break cases_select // 此时中断的是外层的switch执行，而非break所在的for循环。但由于外层switch的分支被中断，内层的for循环也被中断。
				} else {
					println("k  is ", k)
				}
			}
			break // 此break中断的是 switch语句 ，其实不用写这个语句，Go也会在case分支代码都执行完毕后自动执行break，这与C++不同。

		case 50:
			println("now , i is  50")
			break loop //break loop语句中断了外层的对变量i的for循环，使得外层循环不再执行。
		case 70:
			println("now , i is  70")
		}

	}
}
