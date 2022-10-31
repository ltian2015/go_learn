package basic

import (
	"sync"
	"testing"
	"time"
)

/**
  channel(信道)最主要的用途之一就是进程内，多个goroutine连接形成管道式（pipeline）的流计算。
  以channel（信道）作为输入和输出的goroutine可以连接在一起，形成一个管道pipeline，
  每个以channel作为输入和输出goroutine被成为管道的“管段（segment）”
  将各个 “管段”的拼装组成“管道”则是另一个goroutine。
**/

/**
以下是一个使用channel作为函数的输入输出，形成管道（piplelin）

**/
// 这是“源管段（source segment）”主要将接收到的整数数据发送出去。
func gen(nums ...int) <-chan int {

	out := make(chan int)
	sendNums := func() {
		for _, num := range nums {
			out <- num
		}
		close(out)
	}
	go sendNums()
	return out
}

// 这是第二个“管段”，接收第一个“管段”生成的数据，进行平方计算，然后再发送出去。
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	sqThenSend := func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}
	go sqThenSend()
	return out
}

// 这是第最后的“归宿管段（sink segmentn）”，负责打印，不在产生新数据。
func printNum(in <-chan int) {
	print := func() {
		for i := range in {
			println(i)
		}
	}
	go print()

}

// 将三个“管段”组装一个管道，由于channel无缓冲，所以三个管段实际上是串行执行的。
func TestSetupPipeline1(t *testing.T) {
	cout := gen(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)
	sqOut := sq(cout)
	printNum(sqOut)
}

// 下面增加一个“新的管段”可以完成多个“管段”的结果合并，
// 以多个管段作为输入的现象称为“扇入Fan-in”
func merge(cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	sendResult := func(in <-chan int) {
		defer wg.Done()
		for i := range in {
			out <- i
		}
	}
	//为一个输入管道开启一个goroutine去读取数据，然后发送到名为out的统一输出channel（信道）中。
	for _, chIn := range cs {
		go sendResult(chIn)
	}
	//等待所有的goroutine完成工作后，关闭管道，使得读取“管段”可以不在阻塞等待数据的到来。
	go func() {
		wg.Wait()
		close(out)
	}()
	return out //返回工作中的channel给调用者使用
}

// 组装上面新增加merge管段，形成新的管道，合并多个sq的结果。
func TestSetupPipeline2(t *testing.T) {
	ch1 := gen(1, 3, 5, 7, 9, 11, 13, 15, 17, 19)
	ch2 := gen(2, 4, 6, 8, 10, 12, 14, 16, 18, 20)
	sq1 := sq(ch1)
	sq2 := sq(ch2)
	results := merge(sq1, sq2)
	printNum(results)
}

func TestSetupPipeline3(t *testing.T) {
	ch1 := gen(1, 3, 5, 7, 9, 11, 13, 15, 17, 19)
	ch2 := gen(2, 4, 6, 8, 10, 12, 14, 16, 18, 20)
	sq1 := sq(sq(ch1))
	sq2 := sq(sq(ch2))
	results := merge(sq1, sq2)
	printNum(results)
}

/////////////////////////////////////////////////////////////////////
/**组装管道的函数，也是主控函数，如果主控函数想要通过停止上游“源管段”来终止整个“管道”，
该如何实现呢？所有的“管段”函数必须接收一个参数，作为停止信号，而且整个停止信号
必须能够与管道操作“并行”作用，不受管道操作会阻塞后续代码的影响。
做到这一点，就需要将另一个作为信号通知并行channel配合select操作完成。
以下，重构了上述“管段”函数，使得主控函数可以控制管道的停止。
**/
/////////////////////////////////////////////////////////////////////
// 这是“源管段（source segment）”主要将接收到的整数数据发送出去。
func genCancelable(done <-chan int, nums ...int) <-chan int {
	out := make(chan int)
	//需要通过go关键字指定并行执行的函数fp一定不要定义为独立的头等函数，
	//fp函数必须定义在被并行执行它的外层函数f中。
	//如果将fp定义为独立的头等函数，其他人无法知道fp函数是不是需要并行执行。
	// 决定并行执行哪部分代码是函数f的内部细节。因此，f不要对外暴露该细节。
	//除此之外，f作为并行执行了某部分代码的函数，其特征一定是以chan作为输入或（和）输出。
	//看到这种特征的的函数，就应认为函数内部使用的并行技术。
	sendNums := func() {
		defer close(out)
		for _, num := range nums {
			select {
			case <-done: // 正常情况下dong信道会一直阻塞。如果不再阻塞，其已经关闭，表示取消信号到来。
				return
			case out <- num:
				time.Sleep(1000 * time.Millisecond) //为了测试，每秒写一次数据
			}
		}
	}
	go sendNums()
	return out
}

func sqCancelable(done <-chan int, in <-chan int) <-chan int {
	out := make(chan int)
	sqThenSend := func() {
		defer close(out)
		for n := range in {
			select {
			case <-done: // 正常情况下dong信道会一直阻塞。如果不再阻塞，其已经关闭，表示取消信号到来。
				return
			case out <- n * n:
			}
		}
	}
	go sqThenSend()
	return out
}

func printNumCancelable(done <-chan int, in <-chan int) {
	print := func() {
		for i := range in {
			select {
			case <-done:
				return
			default:
				println(i)
			}
		}
	}
	go print()

}

func mergeCancelable(done <-chan int, cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	sendResult := func(in <-chan int) {
		defer wg.Done()
		for i := range in {
			select {
			case <-done:
				return
			case out <- i:
			}

		}
	}
	//为一个输入管道开启一个goroutine去读取数据，然后发送到名为out的统一输出channel（信道）中。
	for _, chIn := range cs {
		go sendResult(chIn)
	}
	//等待所有的goroutine完成工作后，关闭管道，使得读取“管段”可以不在阻塞等待数据的到来。
	go func() {
		wg.Wait()
		close(out)
	}()
	return out //返回工作中的channel给调用者使用
}

func TestSetupPipelineAndCancel(t *testing.T) {
	done := make(chan int)
	ch1 := genCancelable(done, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25)
	ch2 := genCancelable(done, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26)
	sq1 := sqCancelable(done, sqCancelable(done, ch1))
	sq2 := sqCancelable(done, sqCancelable(done, ch2))
	results := mergeCancelable(done, sq1, sq2)
	printNumCancelable(done, results)
	time.Sleep(5000 * time.Millisecond) //等待5秒钟
	close(done)
}
