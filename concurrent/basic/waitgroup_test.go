package basic

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

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

// MD5AllSingleThread函数遍历给定的根目录，提取目录下每个正常文件（不含子目录）的md5数字指纹。
// 并返回文件路径与md5数字指纹字节数组之间的映射（map）
func MD5AllSingleThread(root string) (map[string][md5.Size]byte, error) {
	m := make(map[string][md5.Size]byte)
	fileHandler := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.Type().IsRegular() {
			return nil
		}
		//go1.16版以后读取给定路径文件内容用os.ReadFile函数
		//以前GO 版本则是ioutil.ReadFile函数直接实现该功能，GO1.16之后,
		//ioutil.ReadFile(path)的实现改为调用os.ReadFile，二者效果一样。
		//均比老版本性能高。
		//_, _ := ioutil.ReadFile(path)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		//md5.Sum函数提取文件内容的128位指纹，共计16个字符。
		m[path] = md5.Sum(data)
		return nil
	}
	//filepath.WalkDir遍历一个目录的所有文件，但如果文件是link文件，
	//就不继续访问link文件所指的真正文件或目录
	//filepath.WalkDir()是go1.16之后的新版本，以前使用filepath.Walk()
	//filepath.Walk()相对于WalkDir()的性能低，相对于WalkDir采取了新的策略，
	//包括只向递归函数中传递根文件info的指针而非数据，使用性能更好的fs.DirEntry而非fs.FileInfo。
	//fs.DirEntry数据项比fs.FileInfo更少，但是可以在需要FileInfo内容的时候访问到fs.FileInfo
	err := filepath.WalkDir(root, fileHandler)
	if err != nil {
		return nil, err
	}
	return m, nil
}

/////////////////////////////////////////////////////////////////////////////////////
/**
 下面采用多线程，管道化的并行计算，希望提高程序的执行效率。
 对遍历到的每个文件都开启一个goroutine进行指纹提取。并将结果发送出来。
 再用一个goroutine负责合并所有的结果。
**/
type md5Result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

func MD5AllMultiThreadNoBound(root string) (map[string][md5.Size]byte, error) {
	done := make(chan struct{})
	rc, errc := sumFilesNoBound(done, root)
	select {
	case err := <-errc:
		fmt.Printf("error is : %v\n", err.Error())
		return nil, err
	default:
		return collectSum(rc)
	}

}

// sumFiles函数是一个“管段”函数，管段函数内部执行往往采用多线程并发的模式，
// 而其输入和输出则是“信道(channel)”。多线程并行执行的代码主要负责读写“信道”。
// sumFiles函数开启一个线程遍历每个文件（这样不必等待遍历结束才返回），
// 在该遍历线程中，针对每个遍历到的文件，又开启一个文件处理线程读取该文件内容并提取指纹。
// 因此，sumFiles函数建立并输出了了“一个唯一”的结果信道和一个错误输出信道，
// 每个文件处理线程都将结果写入到结果信道中，如果遍历本身出现错误，则写入错误信道。
// sumFiles作为信道的创建者，要负责结果信道的关闭。
// 所以遍历线程中要创建等待组（sync.WaitGroup）,所有文件处理线程都是等待组中的工作者，
// 当所有工作者都完成工作后（等待组的等待数量为0），则意味着遍历结束，关闭结果通道。
// 这样循环守候处理结果信道的线程就可以不在阻塞，向下执行，从而结束整个处理。
func sumFilesNoBound(done <-chan struct{}, root string) (<-chan md5Result, <-chan error) {
	rc := make(chan md5Result, 5)
	errc := make(chan error)
	//////////////////////////////////
	go func() {
		var wg sync.WaitGroup
		//----------------处理每个被遍历的文件的函数
		fileHandler := func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.Type().IsRegular() {
				return nil
			}
			wg.Add(1)
			//开启一个线程，如果没收到结束信号，则执行文件I/O操作，访问文件数据，并提取指纹。
			//如果收到了结束信号，则什么都不做返回。
			go func() {
				defer wg.Done()
				//执行文件I/O操作，访问文件数据
				data, err := os.ReadFile(path)
				select {
				//提取指纹，如果存在错误，则保存错误。
				case rc <- md5Result{path, md5.Sum(data), err}:
				case <-done:
				}
			}()

			select {
			case <-done:
				return errors.New("遍历操作被取消！")
			default:
				return nil
			}
		}
		//---------------------fileHandler end--------------
		err := filepath.WalkDir(root, fileHandler)
		go func() {
			wg.Wait()
			close(rc)
		}()
		if err != nil {
			errc <- err
		}
	}()
	//////////////////////////////
	return rc, errc
}

/*
*
这个“管段”负责收集个线程向管道中传送的结果，整理成最终结果。
*/
func collectSum(rc <-chan md5Result) (map[string][md5.Size]byte, error) {
	result := make(map[string][md5.Size]byte)
	for md5r := range rc {
		if md5r.err == nil {
			result[md5r.path] = md5r.sum
		}
	}
	return result, nil
}

/////////////////////////////////////////////////////////////////////////////////////

func MD5AllMultiThreadWithBound(root string) (map[string][md5.Size]byte, error) {
	const BOUND int = 8
	done := make(chan struct{})
	//walkFiles开启一个独立线程遍历根目录下所有正常文件路径，输出到pathsc信道中。
	pathsc, errsc := walkFiles(done, root)
	var sumc <-chan md5Result
	select {
	case err := <-errsc:
		fmt.Printf("error is : %v\n", err.Error())
		return nil, err
	default:
		sumc = digestFileStreamByBound(done, pathsc, BOUND)
	}
	return collectSum(sumc)
}

func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	pathChan := make(chan string)
	errChan := make(chan error)
	//独立线程遍历目录，并将遍历root下的所有子目录，将其中正常文件的路径发送出去，以便下一“管段进行处理”
	go func() {
		//注意，信道（channel） outPath由发送者（walkFiles函数）创建，必须由发送者关闭。
		//如果发送者忘记关闭，则很可能导致接收者阻塞（因为接收者往往使用for range 循环进行channel读取）。
		defer close(pathChan)
		handleFile := func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.Type().IsRegular() {
				return nil
			}
			pathChan <- path
			return nil
		}
		err := filepath.WalkDir(root, handleFile)
		if err != nil {
			errChan <- err
		}

	}()
	return pathChan, errChan
}

// 用给定的固定边界数量（bound）的线程（goroutine）提取paths信道（channel）传递过来的path
func digestFileStreamByBound(done <-chan struct{}, paths <-chan string, bound int) <-chan md5Result {
	md5ResultChan := make(chan md5Result)
	var wg sync.WaitGroup
	wg.Add(bound)
	//这是一个独立线程，该线程中开启了固定边界数量（bound）的线程读取paths信道，进行文件的数字指纹提取。
	go func() {
		for i := 0; i < bound; i++ {
			go func() {
				defer wg.Done()
				digester(done, paths, md5ResultChan) //读取paths信道，进行数字指纹提取。
			}()
		}
	}()
	//注意，信道（channel）md5ResultChan由发送者（digestFileStreamByBound函数）创建，必须由发送者关闭。
	//如果发送者忘记关闭，则很可能导致接收者阻塞（因为接收者往往使用for range 循环进行channel读取）。
	//这是常常容易忽视而导致所有线程阻塞而死锁的原因。
	go func() {
		wg.Wait()            //该线程会阻塞等待指定数量的工作线程的结束。
		close(md5ResultChan) //当执行此语句时，表明所有工作线程已经完成任务，要关闭输出信道
	}()
	return md5ResultChan
}

// digester函数是一个真正读取信道进行工作的工作者（worker）函数，该工作者方法使用for range语句主动
// 轮询信道（paths channel）,从中取出数据，进行数据指纹提取的处理，
// 这种通过主动读取信道数进行并发的方式，意味着多个读取同一个信道的并发工作者之间是能者多劳，
// 处理所读取到的数据越快的工作者就会从信道中读取更多的数据，这样可以充分利用计算机的CPU资源。
func digester(done <-chan struct{}, paths <-chan string, md5ResultChan chan<- md5Result) {
	for path := range paths {
		data, err := os.ReadFile(path)
		select {
		case <-done:
			return
		case md5ResultChan <- md5Result{path, md5.Sum(data), err}:
		}
	}
}
