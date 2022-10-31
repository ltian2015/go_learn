package basic

import (
	"fmt"
	"sync"
	"testing"
)

// 知识点1:关于channel的方向。channel天生就是用于goroutines之间传递数据的，因此，数据传递方向天生是双向的（一侧写，一侧读s）。
// 但是，通常使用channel的函数要么是数据发送者（站在channel的写入侧），要么是数据的读取者（站在channel的读取侧）。
// 因此，可以为了确保“单侧”的逻辑可靠性，可以将channle 变量声明名为单向(只读、只写)。
// 而且，由于单方向读或写的channel是不同的类型，可以看作分别实现了read或write接口函数的两个不同类型。
// 而双向channel则可以看作同时实现了read，write两个接口函数的类型，
// 所以，双向channel变量可以赋值给单向channel，但方向不同的单向channel间不可以相互赋值，
// 也不可以把单向channel赋值给双向channel。
// 在channel变量的声明与使用中，箭头方向表示数据流动方向，没有箭头表示数据可双向流动
// 比如：
var writeOnlyChan chan<- int //声明一个只写 channel
var readOnlyChan <-chan int  //声明只读 channel
var readWriteChan chan int   //声明读写双向channel

// 知识点2:关于channel的构建及方向。和其他引用类型map，slice一样，channel在使用前，必须make来创建。
// 否则，如果不用make来创建，则channel的零值是nil，这一点也和其他引用类型（Map、Slice）一样。
// 虽然为了程序逻辑的可靠性，channel变量在单侧（读取侧或写入侧）使用时可以声明为单向，
// 但是创建单向channel几乎不可用。所以，channel创建时都是双向的，而在一侧使用时声明为单向。
// 比如：
func init() {
	if readWriteChan == nil {
		println("readWriteChan is nil ,should be create!")
	}
	readWriteChan = make(chan int) //channel是一个包含了指针的结构，也就是“引用“对象。
	writeOnlyChan = readWriteChan  //channel对象之间的拷贝，只是将”真实对象“的指针地址进行了拷贝。
	readOnlyChan = readWriteChan   //所有拷贝的引用都指向同一个”真实对象“
}

// 知识点3: channel提供的goroutine的基本同步机制——的读写阻塞功能与。完成不同goroutine之间
// 通信功能channel，天生就带阻塞功能，这一功能是其实现goroutine同步机制的基础。具体是：
// 如果使用无缓冲的channel进行通信，因为没有中间的缓存，写入方一旦试图写入(var->chan)，
// 就必须等待接受者准备好，如果接收者没有准备好，那么写入者就会阻塞，无法执行后续操作。
// 而接收方一旦试图接受数据（var<-chan），就会阻塞，直到接收到了写入者所写入的数据。
// 这就好比双人抛接球训练中，传球者必须做好抛球准备（执行向chanenl试图写入操作 var->chan），
// 也就是持球在手等待接收者的准备好接球后,才能将球抛出，球抛出后才能进行下一个动作（解除阻塞）。
// 而接球者一旦试图接球，准备好了接球动作(执行从channel试图取数操作：var<-chan)，
// 就不能有其他动作，直到接到球后才能执行下一个动作(解除阻塞)。
// 无缓冲channel使得两个goroutine遵照“发送一个，处理一个”这种严格的“同步”顺序。
// 如果channel中有缓冲区，则可以同时有多个goroutine对channel进行异步操作，当缓冲区
// 为空时，试图读取channel的“读取侧”goroutine会阻塞，直到缓冲区有数据可以取走。
// 当缓冲区写满时，试图写入channel的”写入侧“goroutine会阻塞，直到缓冲区腾出空间可以写入数据。
// 基于上述认识，在同一个goroutine中对同一个没有缓冲区的channel进行读写“很容易”造成死锁。
// 比如下面两个deadLock()函数：
func TestDeadLock1(t *testing.T) {
	//writeOnlyChan实际上与readOnlyChan是一个实例，所以在同一个goroutine中读写该channel就会死锁
	writeOnlyChan <- 10   //写入数据，此时，当前goroutine会阻塞在此处，下一行无法执行，因为读取侧没有专备好。
	num := <-readOnlyChan //读取数据，但是由于当前goroutine已经阻塞到上一个操作，这个操作不会被执行，又因为该操作是读取侧操作，因此导致死锁。
	println(num)
}

func TestDeadLock2(t *testing.T) {
	//writeOnlyChan实际上与readOnlyChan是一个实例，所以在同一个goroutine中读写该channel就会死锁
	writeOnlyChan <- 10 //写入数据，此时，当前goroutine会阻塞在此行，下一行无法执行，因为读取侧没有准备好。
	go func() {         //读取数据，但是由于当前goroutine已经阻塞在上一个操作上，这个操作不会被执行，又因为该操作是读取侧操作，因此导致死锁。
		num := <-readOnlyChan
		println(num)
	}()
}

// 为了修正死锁问题，必须让两个不同的goroutine读取同一个channel实例。
func TestFixDeadLock(t *testing.T) {
	var wg sync.WaitGroup
	write := func() {
		writeOnlyChan <- 10 //write goroutine试图写入数据，阻塞等待read goroutine准备就绪就发送数据
		wg.Done()
	}
	read := func() {
		num := <-readOnlyChan //read goroutine准备就绪，可以接收写入侧的数据写入。
		println(num)
		wg.Done()
	}
	wg.Add(2)
	go write() //做好发送数据（发球）的准备
	go read()  //做好接收数据（接球）的准备
	wg.Wait()  //等待子例程执行完毕，主例程才退出，否则看不到子例程的执行结果。
}

// 为了修正死锁问题，可以采用有缓冲的channel。
func TestFixDeadLock2(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 10
	num := <-ch
	println(num)
}
func TestWithBufferedChannel(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 10     //发送方（当前gotoutine）只发送了一个数据
	num1 := <-ch // 接收方（当前gotoutine）试图取走第一个数据，因为有缓冲且能成功取走数据，所以不会阻塞。
	num2 := <-ch //接收方（当前gotoutine）试图取走第一个数据，因为没有数据，无法取走数据，所以阻塞，导致程序死锁。
	println(num1)
	println(num2)
}

/**知识点4:关于channel（信道）的关闭。
发送者可通过 close 函数关闭一个信道来“表示没有需要发送的值了”。
必须注意，channel（信道）只能由发送方关闭，接收方不能关闭channel。
向被关闭的channel（信道）发送数据会导致恐慌（panic）异常。
GO没有提供类似isClose()这样判断channel是否关闭的函数，但是，
接收者可以通过为接收表达式分配第二个参数来测试信道是否被关闭：
若没有值可以接收且信道已被关闭，那么在执行完
v, ok := <-ch
之后 ok 会被设置为 false。
读取关闭的通道不会导致恐慌，也不会阻塞，只不过读出的是channel所传递数据类型的空值，
并且，信道与文件不同，通常情况下无需关闭它们。只有在必须告诉接收者不再有需要发送的值时才有必要关闭，
例如终止一个range 循环。
**/

func TestCloseChannel(t *testing.T) {

	ch := make(chan int, 10) //channel缓冲区容量为0
	ch <- 100000             //向channel发送第一个数据
	ch <- 200000             //向channel发送第一个数据
	close(ch)                //关闭channel，关闭带有缓冲区的channel之后，还能否从中读出数据呢？

	i1, isOpen1 := <-ch // 试图从关闭的channel取走被发送的第一个数据
	i2, isOpen2 := <-ch //试图从关闭的channel取走被发送的第二个数据
	i3, isOpen3 := <-ch //试图从关闭的channel取走被发送的第三个数据（实际上并未发送第三个数据）
	//程序实际表明，channel被关闭后，如果带有缓冲区，而且缓冲区中的数据还没有被取走，
	//读取侧还是能够从中取走数据的，而且，试图取没有发送的数据也不会出现异常，只不过取到空值，
	//并且，此时所得到channel的开放（open）标志为false
	println(i1, isOpen1, i2, isOpen2, i3, isOpen3)

	nch := make(chan int)
	close(nch)
	i4, isOpen4 := <-nch //不必等待发送发是否准备好，读取一个已关闭的无缓冲channel不会阻塞，也不会死锁。
	println(i4, isOpen4) //读取的是channel 传递出数据类型的空值，这里是int类型的空值0，isOpen4为false。
}

/**
知识点5: channel的for range 循环读取。
循环 for i := range c 会不断从信道取出数据，直到它被关闭。
如果channel中没有数据，又没有关闭，那么for range 所在goroutine就会阻塞在
for range 这一行（操作）代码上。
**/

func TestForRangeOnChannel(t *testing.T) {
	ch := make(chan int, 5)
	var wg sync.WaitGroup
	send := func(count int) {
		defer wg.Done()
		for i := 0; i < count; i++ {
			ch <- i
		}
		close(ch) //如果不关闭ch,则accept 例程会阻塞。

	}
	accept := func() {
		defer wg.Done()
		for i := range ch {
			fmt.Println("read integer ", i)
		}
	}
	wg.Add(2)
	go send(100)
	go accept()
	wg.Wait() // 等待子例程的结束，否则主例程在子例程之前结束。
}
