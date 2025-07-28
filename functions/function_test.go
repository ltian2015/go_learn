package functions

//这个文件写了一个拱猪游戏的实现，用来展示函数式编程的特点。
import (
	"fmt"
	"math/rand"
	"testing"
	"unsafe"
)

/****************************************************************************
知识点1:什么是函数类型、函数值与函数原型？
	!!! 从GO语言的类型系统的角度看，函数也是GO语言所支持的“类型”之一，与其他类型地位完全一样。
	因此，函数有类型，也有值，完全可以像使用其他类型一样，将函数作为另外一个函数的参数或者返回值。
	!!!函数类型：是一种“组合类型”，也就是“字面类型”。具体而言 ，
		它是由 func 关键字 + 字面的函数签名（ function signature literal）  组合而成，
		而字面的函数签名（ function signature literal）则是由 “输入参数类型列表” 与“结果类型列表”组合而成。
		在实际运用上，可以将func 关键字 + 字面的函数签名看作是函数签名，因此，函数类型就是函数签名。
	!!! 函数原型（function prototype）：是指具体函数声明中除去函数体的部分。也就是 func关键字+ 函数名+ 函数签名。

	!!! 函数值：是程序中，函数签名（函数类型）加上函数体的字面量或字面值。 当我们声明一个自定义函数的时候，
		实际上也声明了一个“不可变的函数值（函数值常量）”。 这个函数值的标识符是函数名，其类型是函数原型去掉函数名的字面类型。
		需要注意的是，系统内建函数不能当作函数值。init函数也不能当作函数值。

	上述概念的范例：
		func (int,string) string就是一个函数签名，也是一个函数类型（字面类型），
		var f func( a int,s string) string 定义一个函数值变量f。
		func f(a int,s string) string { .....} 声明了一个不可变的函数值常量f。

		理解了函数类型与其他类型一样时，我们可以将函数类型做其他函数的输入参数或者结果类型，这就是所谓的高阶函数。
		函数变量实际上是一个指针，指向了函数的代码入口地址，因此，无论函数变量的值是否位nil，64位系统中都会
		占用8个字节的空间。只不过nil值的函数变量没有指向任何函数代码的入口地址，所以调用值为nil的函数会抛出异常。

*******************************************************************************************/
//-----------函数类型与函数值 begin-----------------------------------------//
//下面f1到f5是“定义”的5个函数变量，这些类型完全一样都是func(int, string, string) (int, int, bool)
//而f6则是使用“函数声明”所定义的一个函数变量。变量声明与定义还是不一样，变量声明的同时会赋值，而定义则缺省赋予零值。
//因此，f1到f5这5个函数变量的值是nil，f6不是nil。
var f1 func(int, string, string) (int, int, bool) //标准函数类型方式，
var f2 func(a int, b string, c string) (int, int, bool)
var f3 func(x int, _ string, z string) (int, int, bool)
var f4 func(int, string, string) (x int, y int, z bool)
var f5 func(int, string, string) (a int, b int, _ bool)

// func f6(a int, b string, c string) (int, int, bool)则是一个函数原型(prototype)
// 函数f6是一个函数值常量
func f6(a int, b string, c string) (int, int, bool) {
	var bLen = len(b)
	var cLen = len(c)
	println(a, b, c, bLen, cLen)
	return len(b), len(c), a > bLen+cLen
}

var f7 = f6 //声明的普通函数（函数值常量）可以作为函数值而参与函数变量的赋值

//var f8 = init// 系统自动调用的init函数不能作为函数值
//var f9 = len //系统定义的函数不能作为函数值

func TestFuncTypeAndVariable(t *testing.T) {
	//f1(5, "hi", "lantian")  //空(nil)函数调用会产生panic
	if f1 == nil {
		println("f1 is nil now")
		println("f1 size  is ", unsafe.Sizeof(f1))
		println("f6 size  is ", unsafe.Sizeof(f6))
	}
	f1 = f6
	f2 = f6
	f3 = f6
	f4 = f6
	f5 = f6
	f6(5, "hi", "lantian")
	f1(5, "hi", "lantian")
	f2(5, "hi", "lantian")
}

//-----------函数类型与函数变量 end-----------------------------------------//
/***********************************************************************************
2.关于函数的声明。
     当我们声明一个自定义函数的时候，实际上也声明了一个“不可变的函数值（函数值常量）”。
	 这个函数值的标识符是函数名，其类型是函数原型去掉函数名的字面类型。
	 需要注意的是，系统内建函数不能当作函数值。init函数也不能当作函数值。
 2.1.函数声明中的输入参数、输出结果参数都可以看作函数体的在顶级范围内预先定义的局部变量（local），
   只不过是，这些“函数预置的顶级范围局部变量”一旦声明就已被函数体使用（调用时的赋值操作就是对其的使用）。
   因此，有两个点注意事项：
   （1）.即使函数体内没有使用这些“函数预置的顶级范围局部变量”，也不会出现编译错误，而普通“局部变量”如果生声明
        但不使用，一定会有编译错误提示。
   （2）.函数体的顶级作用域内不能再声明与这些“函数预置的顶级范围局部变量”同名的变量，但下级作用域可以声明。
 2.2.函数原型概念的提出是为了支持函数的声明与函数的实现分离。可以在go文件中使用函数原型声明一个函数，
	然后函数体在.a文件中用go的汇编语言实现。详见https://go.dev/doc/asm

2.3 一个包中的函数不允许重名，但是init函数除外，或者名字为 _ 的函数除外，init函数虽然可以定义多个，但只能被系统自动调用，
    而名字为 _ 的函数虽然可以定义多个，但永远不会被调用,一般被当作还未想好名字的开发过程中的代码交流使用。

 *****************************************************************************************/
//SumImplAssebmly一个自定义的函数声明，这个函数声明定义了一个“常量函数值”，名为SumImplAssebmly，
//这个函数声明只有“函数原型”部分，其函数体在funciplByAssembly.a文件中用go汇编实现。

//func SumImplAssebmly(x, y int64) int64 //函数体在funciplByAssembly.a文件中用go汇编实现。

// 函数f示例用来表明函数参数以及函数体中其他局部变量的作用域。
func f(x, y int) (sum int) {

	var (
	//x int=2        //与函数参数定义冲突
	// y int=3       //与函数参数定义冲突
	//sum int = 5     //与函数参数定义冲突
	)
	{
		var (
			x   int = 1 //与函数参数定义冲突
			y   int = 2 //与函数参数定义冲突
			sum int = 3 //与函数参println(x,y,sum)
		)
		println(x, y, sum)
	}
	sum = x + y
	return
}

// 下列函数可以重名
func init() {
	println("自动调用的初始化函数1")
}
func init() {
	println("自动调用的初始化函数2")
}
func _() {
	println("还没想好叫什么名，无法被调用！")
}
func _() {
	println("还没想好叫什么名，无法被调用！")
}

/******************************************************************************************
3.函数的调用。
  我们所说的函数调用实际就是函数值的使用或者求值（evaluate），函数调用可以看作为对函数值的操作。
  需要注意以下几点：
  !!! 1. 通常函数调用或求值发生在go程序运行期间，但是unsafe包中的函数求值或者调用发生在编译期间，也就是
  !!! 在编译程过程中，被运行的编译器所调用求值。因此，unsafe包中的函数看作常量，可以用于常量表达式定义之中。
  !!!而，有些内置函数，比如 len 和 cap函数既可能在运行时求值，也可能在编译时求值,取决于函数的
  !!! 输入参数是否可以在编译期求值。

  !!! 2.函数调用时的参数值传递方式是“赋值（value assignments）传递” 。所以函数调用时都是拷贝传值，
        如果想要函数改变数据，就必须传递指针。所以对于指针和包含了指针的引用类型就必须要注意，函数的操作可能会修改指向或引用的真实数据。
*******************************************************************************************/

// 一个普通的函数，在运行时求值。
func getIntSize() uint {
	var i int = 0
	return uint(unsafe.Sizeof(i))
}

// const IntSize = getIntSize() //普通函数求值发生在运行时，所以不能赋值给编译时确定值的常量。
const IntSize = uint(unsafe.Sizeof(0)) //unsafe包中的函数是在编译期间由编译器调用求值，因此，可以赋值给常量。

const StrSize int = len("hello") //内建函数len 由于输入参数是常量，所以在编译器求值，可以用在定义常量的表达式中。
var str string = "hello"

//const Str_Size int = len(str) //内建函数len由于输入参数是变量，所以只能在运行期求值，不可以用在定义常量的表达式中。

/**************************************************************************************
4.函数的返回值。
  调用函数，或对函数求值得到的返回值是对函数体内局部的变量的拷贝。需要注意：
  4.1 普通自定义函数的返回值都可以抛弃，但是除了recover 和 copy之外的系统内建函数的返回值不能抛弃，
  比如，len()函数的返回值不可被抛弃，但是可以赋值给匿名变量 _
  4.2 返回值不能抛弃的函数不可以作为 被推迟的出口函数（defer函数）或者goroutine的入口函数。
  4.3 函数的终结除了正常的return 语句返回结果外，还包括异常总结语句，异常终结主要是panic，以及特殊的goto，if、for、select、switch等控制流语句，
      详见，https://go.dev/ref/spec#Terminating_statements
	  由于异常的总结语句会导致函数不会返回，因此返回值没有意义，可以是任意类型。
****************************************************************************************/
//返回结果可以被抛弃的普通函数，为测试函数结果抛弃而定义。
func doSum(a, b int) (int, bool) {
	println("doSum 被调用了！")
	return a + b, true
}

// 测试函数结果的抛弃情况
func TestFuncResultDiscard(t *testing.T) {
	_, ok := doSum(2, 3) //抛弃第一个返回结果
	_ = ok
	sum, _ := doSum(2, 3) //抛弃第二个返回结果。
	doSum(1, 2)           //抛弃所有的返回结果。
	_ = sum
	defer doSum(1, 4)
	go doSum(5, 6)
	str := "hello world"

	//len(str) //系统函数len 的返回结果不允许抛弃，必须作为某个表达式的一部分。
	// defer len(str) //返回结果不可以抛弃的函数不能推迟的出口函数（defer函数）
	// go len(str) //返回结果不可以抛弃的函数不能作为go例程的入口函数。
	_ = len(str) //不可以抛弃的函数结果可以赋值给匿名变量 _
}

// 下列函数都是采用了异常终结方式，没有正常的返回值，所以返回值可以写成任何类型，都不会有编译错误。
func fa1() int {
	panic("error") //异常终结。
}
func fa2() string { //异常终结。
	panic("error")
}

func fb1() int {
a:
	goto a //无限循环，无法正常终结，只能异常终结。
}
func fb2() string {
a:
	goto a //无限循环，无法正常终结，只能异常终结。
}
func fc1() bool {
	for { //无限循环，无法正常终结，只能异常终结。
	}
}
func fc2() int {
	for { //无限循环，无法正常终结，只能异常终结。
	}
}
func fd1() bool {
	if 1 > 2 { //所有分支都无法正常终结的if语句
	a:
		goto a
	} else {
		panic("error")
	}
}
func fd2() string {
	if 1 > 2 { //所有分支都无法正常终结的if语句
	a:
		goto a
	} else {
		panic("error")
	}
}
func fe1(i int) int {
	switch i { //所有分支都无法正常终结的switch语句
	case 0:
		panic("error")
	default:
		panic("error")
	}
}
func fe2(i int) int {
	switch i { //所有分支都无法正常终结的switch语句
	case 0:
		panic("error")
	default:
		panic("error")
	}
}

/***********************************************************************************************
5.Go语言中函数的一些特点总结。
  （1）可以有多个返回值，并可以对返回值命名。
  （2）由于函数类型与值与其他类型与值的地外相同，所以函数值可以作为其他函数的参数，也可以作为结果，支持高阶函数概念
  （3）函数支持可变数量的输入参数。
***********************************************************************************************/
//下面通过一个拱猪游戏的例子展现GO语言函数的这些特点。
const (
	win            = 100 // The winning score in a game of Pig
	gamesPerSeries = 10  // The number of games per series to simulate
)

// A score includes scores accumulated in previous turns for each player,
// as well as the points scored by the current player in this turn.
type score struct {
	player, opponent, thisTurn int
}

// An action transitions stochastically to a resulting score.

// go函数式编程的特点展示1:
// 下面的函数还有一个特点就是对返回的结果参数进行了命名。这里的result和turnIsOver就是命名的返回结果。
// 同时，g函数还可以一个函数可以返回多个值
type action func(current score) (result score, turnIsOver bool)

// roll returns the (result, turnIsOver) outcome of simulating a die roll.
// If the roll value is 1, then thisTurn score is abandoned, and the players'
// roles swap.  Otherwise, the roll value is added to thisTurn.
// go函数式编程的特点展示2:一个函数可以返回多个值
// roll 函数是action类型的一个实例。可以赋值给action类型的变量
func roll(s score) (score, bool) {
	outcome := rand.Intn(6) + 1 // A random int in [1, 6]
	if outcome == 1 {
		return score{s.opponent, s.player, 0}, true
	}
	return score{s.player, s.opponent, outcome + s.thisTurn}, false
}

// stay returns the (result, turnIsOver) outcome of staying.
// thisTurn score is added to the player's score, and the players' roles swap.
// stay 函数也是action类型的一个实例。可以赋值给action类型的变量
func stay(s score) (score, bool) {
	return score{s.opponent, s.player + s.thisTurn, 0}, true
}

// A strategy chooses an action for any given score.

// 函数式编程特点展示3: 高阶函数，一个函数可以输出另一个函数。
type strategy func(score) action

// stayAtK returns a strategy that rolls until thisTurn is at least k, then stays.
func stayAtK(k int) strategy {
	//这里使用了函数字面量来定义一个新函数值，并返回函数值给调用者。
	//当然,也可以把这个函数字面量赋值给父函数中的一个变量，通过这个函数变量进行函数的调用
	//那么就相当于在父函数定义了子函数，虽然GO语言不支持在函数中直接声明子函数，
	//但是采用这种方法也相当于声明了子函数。
	return func(s score) action {
		if s.thisTurn >= k {
			return stay
		}
		return roll
	}
}

// play simulates a Pig game and returns the winner (0 or 1).

// 函数式编程特点展示4: 把具体的函数（值）当作参数来传递。
func play(strategy0, strategy1 strategy) int {
	strategies := []strategy{strategy0, strategy1}
	var s score
	var turnIsOver bool
	currentPlayer := rand.Intn(2) // Randomly decide who plays first
	for s.player+s.thisTurn < win {
		action := strategies[currentPlayer](s)
		s, turnIsOver = action(s)
		if turnIsOver {
			currentPlayer = (currentPlayer + 1) % 2
		}
	}
	return currentPlayer
}

// roundRobin simulates a series of games between every pair of strategies.
func roundRobin(strategies []strategy) ([]int, int) {
	wins := make([]int, len(strategies))
	for i := 0; i < len(strategies); i++ {
		for j := i + 1; j < len(strategies); j++ {
			for k := 0; k < gamesPerSeries; k++ {
				winner := play(strategies[i], strategies[j])
				if winner == 0 {
					wins[i]++
				} else {
					wins[j]++
				}
			}
		}
	}
	gamesPerStrategy := gamesPerSeries * (len(strategies) - 1) // no self play
	return wins, gamesPerStrategy
}

// ratioString takes a list of integer values and returns a string that lists
// each value and its percentage of the sum of all values.
// e.g., ratios(1, 2, 3) = "1/6 (16.7%), 2/6 (33.3%), 3/6 (50.0%)"

// 函数式编程特点展示5: 函数可以接受可变数量的输入参数。
func ratioString(vals ...int) string {
	total := 0
	for _, val := range vals {
		total += val
	}
	s := ""
	for _, val := range vals {
		if s != "" {
			s += ", "
		}
		pct := 100 * float64(val) / float64(total)
		s += fmt.Sprintf("%d/%d (%0.1f%%)", val, total, pct)
	}
	return s
}

func PlayGame() {
	strategies := make([]strategy, win)
	for k := range strategies {
		strategies[k] = stayAtK(k + 1)
	}
	wins, games := roundRobin(strategies)

	for k := range strategies {
		fmt.Printf("Wins, losses staying at k =% 4d: %s\n",
			k+1, ratioString(wins[k], games-wins[k]))
	}
}
