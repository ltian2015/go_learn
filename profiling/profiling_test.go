package profiling

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"testing"
)

/**
  有关性能分析的文章见https://go.dev/blog/pprof
  GO语言进行程序“运行时”的性能分析的两个步骤：
  一、需要在程序中合适的位置写入创建分析报告文件、启动分析、终止分析、关闭分析报告文件的代码。
  二、通过 go tool pprof "pprof file path" 工具打开给路径的分析报告文件进行分析。pprof是一个交互式的命令行工具，
     它有很多子命令，top5,top10等是常用的命令，还有生成各种图形化报告的命令,用go tool pprof help可以查看子命令
  性能分析和提升应该仅针对3%的关键代码。
  图领奖获得者高纳德.克努斯说过：“过早的优化是万恶之源”
  所以，应当在最后根据程序的整体体现，再决定是否进行关键代码的性能优化。
**/
// 这是一个耗费性能的算法程序，是被进行性能分析的对象。
func foo(n int) string {
	var buf bytes.Buffer
	for i := 0; i < 100000; i++ {
		s := strconv.Itoa(n)
		buf.WriteString(s)
	}
	sum := sha256.Sum256(buf.Bytes())
	var b []byte
	for i := 0; i < int(sum[0]); i++ {
		x := sum[(i*7+1)%len(sum)] ^ sum[(i*5+3)%len(sum)]
		c := strings.Repeat("abcdefghijklmnopqrstuvwxyz", 10)[x]
		b = append(b, c)
	}
	return string(b)
}
func TestCpuProfiling(t *testing.T) {
	//!!!创建一个文件用于存储CPU分析内容
	cpufile, err := os.Create("cpu.pprof")
	if err != nil {
		panic(err)
	}
	//!!!启动CPU分析
	err = pprof.StartCPUProfile(cpufile)
	if err != nil {
		panic(err)
	}
	defer cpufile.Close()        //!!!关闭CPU分析文件
	defer pprof.StopCPUProfile() //!!!停止CPU分析—先于关闭CPU分析文件执行

	//!!!! 正式运行被测试的代码
	if foo(12345) == "aajmtxaattdzsxnukawxwhmfotnm" {
		fmt.Println("Test PASS")
	} else {
		fmt.Println("Test FAIL")
	}
	for i := 0; i < 100; i++ {
		foo(rand.Int())
	}
}

func cpuProfiling(fn func(), profilingFilePath string) {
	cpufile, err := os.Create(profilingFilePath)
	if err != nil {
		panic(err)
	}
	//启动CPU分析
	err = pprof.StartCPUProfile(cpufile)
	if err != nil {
		panic(err)
	}

	fn()
	defer cpufile.Close()        //关闭CPU分析文件
	defer pprof.StopCPUProfile() //停止CPU分析—先于关闭CPU分析文件执行
}
func calcl() {
	//确保程序输出是正确的
	if foo(12345) == "aajmtxaattdzsxnukawxwhmfotnm" {
		fmt.Println("Test PASS")
	} else {
		fmt.Println("Test FAIL")
	}
	for i := 0; i < 100; i++ {
		foo(rand.Int())
	}
}
func TestCpuProfiling2(t *testing.T) {

	cpuProfiling(calcl, "cpu2.pprof")
}
