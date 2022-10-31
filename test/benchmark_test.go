/*
*

进行基准测试的文件必须以*_test.go的文件为结尾，这个和测试文件的名称后缀是一样的，例如abc_test.go
参与Benchmark基准性能测试的方法必须以Benchmark为前缀，例如BenchmarkABC()
参与基准测试函数必须接受一个指向Benchmark类型的指针作为唯一参数，*testing.B
基准测试函数不能有返回值.

b.ResetTimer是重置计时器，调用时表示重新开始计时，可以忽略测试函数中的一些准备工作
b.N是基准测试框架提供的，表示循环的次数，因为需要反复调用测试的代码，才可以评估性能。
因而，普通的基准测试函数如下 ：
func BenchmarkXXXXXXX(b *testing.B){

	for i := 0; i < b.N; i++ {
	        //!!! 被测试的代码或函数
	}

}
!!!  如果被测用例需要多个可变的输入参数，往往采用表驱动测试模式，在这种模式下，
!!!  要使用b.Run() 方法来驱动类似上面的测试。具体情况，可见下面面具体案例。

go test 命令执行测试，进入测试用例所在包所在的目录，然后执行go  test 命令；
下面列出了go test 命令的使用方法:
参数            		含义
-bench=regexp	性能测试，支持表达式对测试函数进行筛选。
-bench=.		则是对当前目录（包）下所有的benchmark函数测试，指定名称则只执行具体测试方法而不是全部
-benchmem		性能测试的时候显示测试函数的内存分配的统计信息
－count n		运行基准测试的次数，默认一次。
-benchtime		Ns基准测试运行的时间（几轮测试不定），也可指定一轮测试中运行的次数（Nx）
-run regexp		只运行特定的测试函数， 比如-run ABC只测试函数名中包含ABC的测试函数. ^# 可以屏蔽非基准测试的其他测试函数
-timeout t		测试时间如果超过t, panic,默认10分钟
-cpuprofile=fileName	以二进制格式存储被测程序中调用的函数的CPU使用详情，用 go tool pprof fileName可以查看该文件各项内容。
-memprofile=fileNamme	以二进制格式存储被测程序中调用的函数的内存使用详情，用 go tool pprof fileName可以查看该文件各项内容。
-v				显示测试的详细信息，也会把Log、Logf方法的日志显示出
!!!  在进行基准测试时会运行所有非基准测试的测试函数，使用-run=^#可以阻止非基准测试之外的其他测试函数的运行，
!!!  如下：

	go test   -bench='GunzipNopool$'  -run=^#   -benchmem

!!! 比较好的资料如下：

	https://blog.logrocket.com/benchmarking-golang-improve-function-performance/

!!! 使用重定向操作可以将测试结果输出到文件，但使用 管道命令 tee 不仅可以将测试结果输出到文件，还能在屏幕上显示。
比如：

	go test   -bench='BenchmarkPrimeNumbersDrivenByTable'  -run=^#   -benchmem  -count=5 | tee testReport.txt

!!! 使用 benchstat BenchReprotFile 命令可以对基准测试结果（文件）进行统计。
benchstat需要下载，下载命令为 ：
go install golang.org/x/perf/cmd/benchstat@latest
对 一个测试结果中的几次测试进行对比分析：
$ benchstat testReport.txt
对两个测试结果中的几次测试进行对比分析：
$ benchstat testReport.txt testReport_new.txt

!!! 如果想更深层次得知基准测试中，被测程序中性能开销较大的函数数调用的情况，需要
!!! 使用参数 -cpuprofile=filename、参数-memprofile=fileNamme、-blockprofile=fileName,将会
!!! 生成CPU、内存和阻塞的分析文件，再配合 go tool pprof  filenamme工具查看具体情况。
!!! 启动 CPU 分析时，运行时(runtime) 将每隔 10ms 中断一次，记录此时正在运行的协程(goroutines)
!!! 的堆栈信息。程序运行结束后，可以分析记录的数据找到最热代码路径(hottest code paths)。
!!! 一个函数在性能分析数据中出现的次数越多，说明执行该函数的代码路径(code path)花费的时间占总运行时间的比重越大。
!!!  内存性能分析(Memory profiling) 记录堆内存分配时的堆栈信息，忽略栈内存分配信息。
!!! 内存性能分析启用时，默认每1000次采样1次，这个比例是可以调整的。因为内存性能分析是基于采样的，
!!! 因此基于内存分析数据来判断程序所有的内存使用情况是很困难的。
!!! 阻塞性能分析(block profiling) 是 Go 特有的。
!!! 阻塞性能分析用来记录一个协程等待一个共享资源花费的时间。在判断程序的并发瓶颈时会很有用。阻塞的场景包括：
!!! 在没有缓冲区的信道上发送或接收数据。
!!! 从空的信道上接收数据，或发送数据到满的信道上。
!!! 尝试获得一个已经被其他协程锁住的排它锁。
!!! 一般情况下，当所有的 CPU 和内存瓶颈解决后，才会考虑这一类分析。

启用CPU与内存分析的基准测试范例如下：
go test   -bench='BenchmarkPrimeNumbersDrivenByTable'  -run=^#   -benchmem  -cpuprofile=cpu.out -memprofile=mem.out
再用
go tool pprof cpu.out查看所涉及的各个函数的CPU使用情况，
go tool pprof mem.out查看所涉及的各个函数的内存使用情况，

pprof 中的top子命令可以查看top10 CPU开销的函数,top命令返回占用资源前10的函数调用。
以CPU 资源为例，top子命令返回结果如下：
flat：指的是该方法所占用的资源（这里是CPU时间）（不包含这个方法中调用其他方法所占用的时间）
flat%: 指的是该方法flat资源（这里是时间）占全部采样时间的比例。
cum：指的是该方法以及方法中调用其他方法所占用资源（这是CPU时间）总和，这里注意区别于flat
cum%:指的是该方法cum资源（这里是CPU时间）占全部采样资源的比例
sum%: 指的是执行到当前方法累积占用的CPU时间总和，也即是前面flat%总和
!!!  flat%是一个非常重要的指标，表示了一个函数自身在整个采样时间内的资源占用。
可以用peek regexp 查看名称符合regexp正则表达式的函数的内存的峰值使用。
可以用help子命令查看pprof所有用法。

!!! 当然，可以通过WEB界面可视化的方式查看分析结果，比如：
$ go tool pprof -http=:9999 cpu.pprof
!!! 但，这需要安装Graphviz 工具 ，
如果提示 Graphviz 没有安装，则通过 brew install graphviz(MAC) 或 apt install graphviz(Ubuntu) 即可。
*
*/
package test

import (
	"fmt"
	"math"
	"testing"
)

// 待测的业务函数
func primeNumbers(max int) []int {
	var primes []int

	for i := 2; i < max; i++ {
		isPrime := true

		for j := 2; j <= int(math.Sqrt(float64(i))); j++ {
			if i%j == 0 {
				isPrime = false
				break
			}
		}

		if isPrime {
			primes = append(primes, i)
		}
	}

	return primes
}

var num = 1000

// 基本的基准测试函数
func BenchmarkPrimeNumbersBasic(b *testing.B) {
	for i := 0; i < b.N; i++ { //运行N次测试程序以便获得可靠的测试结果
		primeNumbers(num)
	}
}

// !!! 由于基准测试函数只允许输入一个类型为*testing.B的参数，
// !!! 所以要对不同的输入进行测试，就要采用表驱动测试模式（table driven test pattern），
// !!! 然后在基准测试函数中使用 b.Run() 方法
var table = []struct {
	input int
}{
	{input: 100},
	{input: 1000},
	{input: 74382},
	{input: 382399},
}

// 表驱动模式的基准测试函数，使用以下测试命令进行基准测试，并将结果保存到文件中。
// $ go test   -bench='BenchmarkPrimeNumbersDrivenByTable'  -run=^#   -benchmem  -count=5 | tee testReport.txt
func BenchmarkPrimeNumbersDrivenByTable(b *testing.B) {
	for _, v := range table { //遍历表

		b.Run(fmt.Sprintf("input_size_%d", v.input), func(b *testing.B) { //b.Run驱动基准测试
			for i := 0; i < b.N; i++ { //运行N次测试程序以便获得可靠的测试结果
				primeNumbers(v.input)
			}
		})
	}
}
