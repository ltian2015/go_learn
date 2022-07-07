package basic

import (
	"crypto/md5"
	"fmt"
	"testing"
	"time"
)

func TestDeadLockBychannel(t *testing.T) {
	DeadLock1()
	println("hi,hello")
}
func TestDeadLockFix(t *testing.T) {
	FixDeadLock()
	println("hi,hello")
}
func TestDeadLockFix2(t *testing.T) {
	FixDeadLock2()
	println("hi,hello")
}
func TestBufferedChannel(t *testing.T) {
	BufferedChannel()
}
func TestChannelClose(t *testing.T) {
	CloseChannel()
}

func TestForRangeOnChannel(t *testing.T) {
	ForRangeOnChannel()
}

func TestSelectCaseEvaluate(t *testing.T) {
	SelectCaseEvaluate()
}

func TestPipeline1(t *testing.T) {
	SetupPipeline1()
}

func TestPipeline2(t *testing.T) {
	SetupPipeline2()
}

func TestPipeline3(t *testing.T) {
	SetupPipeline3()
}

func TestPipelineCancel(t *testing.T) {
	SetupPipelineAndCancel()
}
func doMd5Test(path string, testedFunc func(path string) (map[string][md5.Size]byte, error)) {
	start := time.Now()
	fdMap, err := testedFunc(path)
	//求取时间差，求取时间差(Duration)方法包括：Time对象的Sub方法，time包中的Since和Until函数
	dur := time.Since(start) //等价与以下两句
	//end := time.Now()
	//durMs := end.Sub(start)
	fmt.Printf("程序持续了 %v毫秒\n", dur.Milliseconds())
	if err != nil {
		println(err)
		return
	}
	for k, v := range fdMap {
		fmt.Printf("%v %x\n", k, v)
	}
	fmt.Printf("共提取了%v个文件的数字指纹\n", len(fdMap))
}
func TestSingleThreadMd5All(t *testing.T) {
	const root string = "/Users/learn/Documents/_development/c"
	doMd5Test(root, MD5AllSingleThread)
}

func TestMD5AllMultiThreadNoBound(t *testing.T) {
	const root string = "/Users/learn/Documents/_development/c"
	doMd5Test(root, MD5AllMultiThreadNoBound)
}

func TestMD5AllMultiThreadByBound(t *testing.T) {
	const root string = "/Users/learn/Documents/_development/c"
	doMd5Test(root, MD5AllMultiThreadWithBound)
}
