/*
*
!!! go中进行debug测试，要求必须安装dlv工具，但dlv工具版本与go版本不匹配时则无法进行测试与debug，
因此，必须更新dlv工具，但更新安装无法自动覆盖，必须要删除原有已安装的dlv工具，更新方式安装dlv的命令如下：
sudo rm $GOROOT/bin/dlv
sudo rm $GOPATH/bin/dlv
sudo go get -u github.com/go-delve/delve/cmd/dlv

go的自动化自测试工具go testd对被测试的文件与函数的命名约定：

	1.测试文件名必须以  _test结尾
	2.测试函数必须以Test开头，参数必须是(t *testing.T)
	3.一个包路径下只有一个测试文件。（在vscode工具至少是这样）

按照上述约定，在vscode中，使用go test -v ${fileDirname}就可以测试当前文件所在包。
-v 参数可以把测试过程中的fmt.print 结果显示出来。

controlflow包主要学习for,if,switch等三个用于控制程序执行分支的语句。
*
*/
package controlflow

/**
本文件演示了for循环的几种用法，以及如何用标签实现嵌套循环的灵活终止控制。
for语句的功能用来指定重复执行的语句块，for语句中的表达式有三种：
ForStmt = "for" [ Condition | ForClause | RangeClause ] Block .
Condition = Expression .
ForClause = [ InitStmt ] “;” [ Condition ] “;” [ PostStmt ] .
RangeClause = [ ExpressionList “=” | IdentifierList “:=” ] “range” Expression .

由语法可以看出，for语句有三种可选的表达法，即：
1. for InitStmt ;  Condition ; PostStmt { block}， 正常的for循环
2. for Condition  { block } 相当于while循环
3. for RangeClause { block}  for与range配合，相当于for each
4. for{ block} 无限循环
并且，当for与 range 配合时，只允许array,slice,string,map ,channel这个5个包含多个同类元素的类型。
Range expression                          1st value          2nd value

array or slice  a  [n]E, *[n]E, or []E    index    i  int    a[i]       E
string          s  string type            index    i  int    see below  rune
map             m  map[K]V                key      k  K      m[k]       V
channel         c  chan E, <-chan E       element  e  E

**/
import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

func TestFor() {
	classicForLoop()
	wihleLoopWithFor()
	infiniteLoop()
	foreachLoop()
	//testForBreak()
	testForLoopLabel()
}

//经典的For循环，
// for 初始化语句;是否执行执行体的判断条件语句;执行体执行后的后置处理语句 { 执行体}

func classicForLoop() {
	var sum int = 0
	for i := 0; i <= 1000; i++ {
		sum = sum + i
	}
	println(sum)
}

// wihleLoopWithFor函数用for 实现while，因为golang中没有while循环。
func wihleLoopWithFor() {
	const start = 1
	const target = 100
	i := start
	//相当于 while(i< target)
	for i < target {
		i = i * 2
		println(i)
	}
}

// infiniteLoop函数使用for实现无限循环。for{ }就是无限循环。
func infiniteLoop() {
	fmt.Print("Please enter a  number :")
	for {
		reader := bufio.NewReader(os.Stdin)
		char, _ := reader.ReadByte()
		if char == '8' {
			println("you succeed!")
			break
		}
		fmt.Print("Fail, try enter other number :")
	}

}

// foreachLoop函数使用for与range操作 配合，实现foreach，
// range 操作目前只适合array,slice,string,map ,channel这个5个包含多个同类元素的类型。
// 必须注意的是，for与range操作配合的时候，会把每个元素及其index拷贝赋值给for块中临时定义的二元组变量，
// 对该变量取地址，会得到同一个地址，试图通过修改该变量的值来修改被遍历的元素都不会有效果。
func foreachLoop() {
	type Foo struct {
		bar string
	}
	list := []Foo{{bar: "A"}, {bar: "B"}, {bar: "C"}}
	cp := make([]*Foo, len(list))
	// 注意，value是for 范围内的一个变量，list中的每个元素都拷贝到这个变量中，
	// 对该变量做取地址操作， 只能取到同一个地址。
	// 试图通过该变量修改所便利的列表中的元素是不可能做到的，除非列表中的元素是指针。
	for i, value := range list {
		cp[i] = &value //得到的都是同一个地址，也就是value的临时地址。
	}
	fmt.Printf("cp:%v\n", cp)

	var s = "hello"                                   //字符串
	var ia [5]int = [5]int{1, 2, 3, 4, 5}             //数组
	var aSlice []int = make([]int, 3)                 //切片
	var amap map[int]string = make(map[int]string, 5) //映射
	var channel1 chan int = make(chan int, 6)
	//对字符串进行range操作，得到的第一个变量是字符在字符串中index，第二个变量是字符
	for i, abyte := range s {
		fmt.Printf("%v : %c\n", i, abyte)
	}
	//对数组进行range操作，得到的第一个变量是数组元素的index，第二个变量是数组元素
	var sum int = 0
	//用空标识符 _ 取代程序中用不到的数组元素的序号index
	for _, e := range ia {
		sum = sum + e
	}
	println(sum)
	//对切片进行数组操作，得到的第一个变量是元素在切片中的index，第二个变量是切片元素
	for i, _ := range aSlice {
		aSlice[i] = ia[i]
	}
	for _, element := range aSlice {
		println(element)
	}
	//对
	for i := 0; i < 5; i++ {
		amap[i] = fmt.Sprintf("str%v", i)
	}
	amap[6] = "str6"
	amap[7] = "str7"
	amap[7] = "str8"
	for k, v := range amap {
		fmt.Printf("%v : %s\n", k, v)
	}
	// 协调3个routine的同步完成，一个向channel中写数据的进程，两个取出数据的线程。
	var wg sync.WaitGroup
	wg.Add(3)
	//运行一个子routine,每隔10毫秒向channel中写入一个数据。
	go func() {
		defer wg.Done() //确保在各种情况下(包括异常)的结束，都会使阻塞锁的数量减少。
		for i := 0; i < 30; i++ {
			time.Sleep(10 * time.Millisecond)
			channel1 <- i
		}
		close(channel1) //关闭channel不会导致没有被取走的元素无法被其他goroutine取走，只是channel变空时不会阻塞读取routine。
	}()
	//定义一个闭包函数，取出channel1中的元素进行打印处理。
	readChannel := func(routineName string) {
		defer wg.Done() //确保在各种情况下(包括异常)的结束，都会使阻塞锁的数量减少。
		//如果channel 中还有元素没有被取出，那么range操作就不会阻塞，否则就会阻塞，直到channel被关闭。
		for e := range channel1 {
			println(routineName, " routine read ", e)
		}
	}
	//运行两个子routine并行取出channel中的元素
	go readChannel("routine1")
	go readChannel("routine2")
	//主routine等待所有子routine处理结束;
	//否则，子routine虽在运行，而主routine会立即结束会导致整个程序结束，就看不到子routine的结果
	wg.Wait() //等待三个子routine的执行完毕！
	println("the end!")
}

// for 循环上面可以定义一个label，使用goto命令可以退出到标记之处，重新来一轮循环
func testForBreak() {
	var loopTimes = 1
	//标号声明只能放在for循环语句前面，二之间不能间隔其他有效语句行
loop:
	for i, j, rd := 1, 1, rand.New(rand.NewSource(time.Now().Unix())); i*j < 10000; rd, i, j = rand.New(rand.NewSource(time.Now().Unix())), rd.Intn(100), rd.Intn(100) {
		time.Sleep(3 * time.Nanosecond)
		if i+j > 100 && i+j < 150 { //出现两个随机数之和在(100,150)之间，则重启一轮循环
			fmt.Printf(" the %v  loop  times is over, start a new loop  with i=%v j=%v \n", loopTimes, i, j)
			loopTimes++
			goto loop // 中断本轮循环，再启动一轮循环。goto 跳转到标签可以重新开启循环
		} else if i+j >= 150 { //两个随机数之和大于等于150，就结束循环
			println("break loop with i=", i, " j=", j, " at the ", loopTimes, " loop times")
			break loop //从本轮循环操作中彻底退出,跳到loop标签处，但不再执行循环，即不进行新一轮循环呢；与break等同。
		}
		if i > 50 { //如果出现i大于50，打印数据然后，进行下一次迭代。
			println("continue with i=", i, " j=", j)
			continue //不再执行下面代码，至今进行新的迭代
		}
		println("i=", i, " j=", j, " i*j=", i*j, " i+j=", i+j)
	}
}

// 在go中，只有 break,continue ,goto 3个命令可以操作标签(label)
// 标签的妙用是：为嵌套循环定义标签(label)有利于在内外循环中控制外层循环的break或continue。
// 注意，尽量不要用goto配合标签来实现是逻辑分支的跳转，在大程序中，这种方式不利于“程序的构块化”。
func testForLoopLabel() {
outerLoop:
	for i := 0; i < 100; i++ {
	inerloop:
		for j := 0; j < 100; j++ {
			if j == 50 {
				println("break inner loop when j==50 ")
				break inerloop //如果不写中断到何处，则默任中断所在循环体的执行。
			}
			if i == 75 && j == 49 {
				println("break outer loop when i==75 and j==49 ")
				break outerLoop // 在内循环中结束外循环
			}
		}
		//实际上执行不到了，因为内循环已经在i==75时终止了内循环。
		if i == 80 {
			println("break outer loop when i==80 ")
			break //如果不写中断到何处，则默任中断所在循环体的执行。
		}
	}
}
