package test

//!!! 这个包主要根据CPU的自身构成及CPU一级缓存的工作原理来提升的几种方法，分别是：
//!!!  1.对密集循环进行循环展开，提高CPU管道中相互不依赖的指令之间的交错执行来减少停顿等待
//!!!  2.函数调用时采用寄存器(register)而不是栈（stack）进行函数的参数与返回值传递。
//!!!  3. 通过对大数据集进行排序来提供CPU的分支预测能力，减少预测不准而导致的进入管道的指令的回退与重新执行。
//!!!  4.利用CPU一级缓存的整行加载的特性提升大数据集的处理效率
import (
	"math/rand"
	"slices"
	"testing"
)

// /////////////////////////////////////////////////////////////////////////////////
// /////////!!! 循环展开///////////////
// 紧循环
func tightLoop(num int) int {
	var result int = 1
	//!!! 紧循环d代码
	for i := 1; i <= num; i++ {
		result += i
	}
	println(result)
	return result
}

// 循环展开
func loopUnrolling(num int) int {
	var result int = 1
	//!!! 循环展开
	for i := 1; i <= num; i = i + 4 {
		result += i
		result += i + 1
		result += i + 2
		result += i + 3
	}
	println(result)
	return result
}

// 循环展开且减少语句间的依赖，使得不依赖的指令与相互依赖的指令可以交错进入管道同时执行
func loopUnrollingWithDependencyReduce(num int) int {
	var result int = 1
	for i := 1; i <= num; i = i + 4 { //!!! 循环展开；循环的步长变大，减少条件判断的次数，降低分支惩罚的代价。
		//!!!该循环体将calcFactorial1中的一条代码写为4条，是为了发挥“管道流水线（Pipelining）”架构下CPU的优势，
		//!!! 这样，CPU可以才采用“变序执行（out-of-order exectuion）”，将更多不依赖的指令加载到Pipeline中同时执行，
		//!!! 减少因依赖所引起CPU流水线阶段（pipeline stage） 的空转
		//注意，将循环体中的一条指令展开为多条的优化方式叫做“循环展开（loop unrolling）”
		//注意，循环体中的代码，除了最后一条外，其他每条都减少了对上一条代码输出结果的依赖。这会提高多核的并发能力。
		sum1 := result + i
		sum2 := i + 1
		sum3 := i + 2
		sum4 := i + 3
		result = sum1 + sum2 + sum3 + sum4
	}
	println(result)
	return result
}

func BenchmarkTightLoop(b *testing.B) {
	tightLoop(1000000000)

}

func BenchmarkLoopUnrollingWithDependencyReduce(b *testing.B) {
	loopUnrollingWithDependencyReduce(1000000000)

}

// /////////////////////////////////////////////////////////////////////////////////
// ///!!! 使用寄存器传递函数调用参数与结果    ////////
//
//	sum函数被calcFactorial3反复大量调用，具有优化的可能
func sumFuncReturnByStack(a, b int64) int64 {
	return a + b
}

func callSumFuncReturnByStack(num int64) int64 {

	var result int64 = 1
	var i int64
	for i = 1; i <= num; i++ {
		result = sumFuncReturnByStack(result, i)
	}
	println(result)
	return result

}

// !!! 如果一个函数会被其函数大量的调用，那么在1.17版本后，go:noinline这个编译指令可以让该函数
// !!! 被编译为采用寄存器而不是堆栈的参数/返回值传递的调用规约，提高性能。
// !!! 优化方方法如下：
//
//go:noinline
func sumFuncReturnByRegister(a, b int64) int64 {
	return a + b
}
func callSumFuncReturnByRegister(num int64) int64 {

	var result int64 = 1
	var i int64
	for i = 1; i <= num; i++ {
		result = sumFuncReturnByRegister(result, i)
	}
	println(result)
	return result

}

func BenchmarkFuncReturnByStack(b *testing.B) {
	callSumFuncReturnByStack(1000000000)

}
func BenchmarkSumFuncReturnByRegister(b *testing.B) {
	callSumFuncReturnByRegister(1000000000)
}

// ///////////////////////////////////////////////////////////////////////////////////////
// !!! 编译器可以对程序中(if之类语句引发的)指令分支进行预测，从而保证CPU处于循环中的指令分支沿着同一个分支前进，不用切换分支，提高性能。
// !!! Pipelining结构的CPU，如果连续的两条指令不在一个分支中，将会导致管道中的指令的浪费，回退和重新执行，详见：
// !!! https://stackoverflow.com/questions/11227809/why-is-processing-a-sorted-array-faster-than-processing-an-unsorted-array
// !!! 为此，影响程序性能的关键代码中，要提高编译器对（if之类语句所引发的）指令分支预测成功率。
func calcUnSortArray() {
	const SIZE = 32768
	var data []int = make([]int, SIZE)
	r := rand.New(rand.NewSource(2))
	for i := 0; i < SIZE; i++ {
		data[i] = r.Int() % 256
	}
	println(data[0], data[1], data[2])
	//!!! 注意，这里未对切片进行排序，使得编译器的分支预测优化器失效。
	//slices.Sort[[]int, int](data)
	var sum int
	for j := 0; j < 100000; j++ {
		for k := 0; k < SIZE; k++ {
			if data[k] >= 128 { //!!!编译器分支预测期对此处分支进行优化,优化失败
				sum += data[k]
			}
		}
	}
	println("sum=", sum)

}

// !!! 切片/数组等大数据集合的排序会提高分支优化性能
func calcSortedArray() {
	const SIZE = 32768
	var data []int = make([]int, SIZE)
	r := rand.New(rand.NewSource(2))
	for i := 0; i < SIZE; i++ {
		data[i] = r.Int() % 256
	}
	println(data[0], data[1], data[2])
	//!!! 注意，这里对切片进行了排序，使得编译器的分支预测优化器得以生效。
	//!!! 如果没有下面这一行，不会影响逻辑，但是程序的性能会下降一个数量级。
	slices.Sort(data)
	var sum int
	for j := 0; j < 100000; j++ {
		for k := 0; k < SIZE; k++ {
			if data[k] >= 128 { //!!!编译器分支预测期对此处分支进行优化,优化成功
				sum += data[k]
			}
		}
	}
	println("sum=", sum)

}
func BenchmarkCalcUnSortArray(b *testing.B) {
	calcUnSortArray()

}
func BenchmarkBranchCalcSortedArray(b *testing.B) {
	calcSortedArray()

}

// !!! 测试if语句的分支预测性能，为了与case语句对比；
// !!!测试结果表明二者似乎没有差别,可能是因为case和if都是高级的语法糖，最终使用相同CPU指令跳转指令
func ifBranch() {
	rn := rand.New(rand.NewSource(2))
	for i := 0; i < 1000000; i++ {
		n := rn.Intn(5) % 5
		if n == 0 {
			print(0)
		} else if n == 1 {
			print(1)
		} else if n == 2 {
			print(2)
		} else if n == 3 {
			print(3)
		} else if n == 4 {
			print(4)
		} else {
			panic("rand error")
		}
	}
}

// !!!测试case分支语句的分支预测性能，为了与上述if分支语句比比较；
// !!!测试结果表明二者似乎没有差别,可能是因为case和if都是高级的语法糖，最终使用相同CPU指令跳转指令
func caseBranch() {
	rn := rand.New(rand.NewSource(2))
	for i := 0; i < 1000000; i++ {
		n := rn.Intn(5) % 5
		switch n {
		case 0:
			print(0)
		case 1:
			print(1)
		case 2:
			print(2)
		case 3:
			print(3)
		case 4:
			print(4)
		default:
			panic("rand error")
		}
	}
}
func BenchmarkCaseBranch(b *testing.B) {
	caseBranch()
}

func BenchmarkIfBranch(b *testing.B) {
	ifBranch()
}

///////////////////////////////////////////////////////////////////////////////////////
///!!! 利用内存缓存的工作原理提高大数据集数据处理的程序性能/////////////////////
