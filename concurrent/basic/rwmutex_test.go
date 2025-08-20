package basic

import (
	"sync"
	"testing"
	"time"
)

/**
!!! RWMutex 是一种读写互斥锁。
!!! 该锁可以由"任意数量"的读取者(goroutine)或"单个"写入者(goroutine)持有,
!!! 也就是说，[RWMutex.RLock] 可以被多个读取者(goroutine)持有，不会阻塞，
!!! 而[RWMutex.Lock] 只能被单个goroutine，往往是作为写入者的goroutine所获得。
!!! RWMutex 的零值表示未锁定的互斥锁。
!!! RWMutex 在首次使用后不得复制。
!!! [RWMutex.RLock] 不能升级为 [RWMutex.Lock]
!!! [RWMutex.Lock] 也不能降级为 [RWMutex.RLock]。

!!! 在多goroutine 并行执行的场景下，如果某个 goroutine 在”读取锁“已被一个或
!!! 多个读取者（goroutine）持有的情况下，调用了[RWMutex.Lock]，那么
!!! !!! 对 [RWMutex.RLock] 的并发调用将阻塞，直到写入者（goroutine）获取（并释放）锁，
!!! 以确保锁最终可供写入者使用。并且，写入者（goroutine）即使在持有了写锁的时候，
!!! 也必须阻塞等待所有读锁解锁才能执行，也就要是等待所有已获得了”读锁”的goroutine
!!! 处理所有的读取操作之后解出读取锁才能执行，如果被写入者（goroutine）等待的读取者（goroutine）
!!! 忘记了释放自己写锁，那么写入者就会一直等待下去（期望他人释放锁），这就会造成“死锁”。

请注意，这禁止了递归的读锁。
用 Go 内存模型的术语来说，当 n < m 时，第 n 次调用 [RWMutex.Unlock] 会“同步于”第 m 次
调用 Lock，就像 [Mutex] 一样。对于任何 RLock 调用，都存在一个 n，
使得第 n 次调用 Unlock 会“同步于”该 RLock 调用，
而相应的 [RWMutex.RUnlock] 调用会“同步于”第 n+1 次调用 Lock。
**/

// TestRwmutexLockRlockSequence测试先Lock后Rlock的顺序关系。
// 注意，测试表明，"写锁"操作(Lock)施加后，其他线程无法再获取”读取锁“操作（RLock），
// !!!这些试图获取读所的线程会阻塞，直到写锁被释放。
// !!! Lock
func TestRwmutexLockRlockSequence(t *testing.T) {
	var _lock = sync.RWMutex{}
	var _count = 0
	var _wg sync.WaitGroup
	// !!!count 变量的读取者，会被30并行的goroutine同步调用
	reader := func(i int) {
		_wg.Add(1)
		defer _wg.Done()
		_lock.RLock()
		defer _lock.RUnlock()
		seeCount := _count
		time.Sleep(2 * time.Second)
		//!!! 读取锁的存在，使得即使等待两秒钟,这段时间也会有其他改写_count的线程在执行，
		//!!! 也不会破坏本线程的数据一致性。就是说使得两秒前后，seeCount和_count是同一个数据
		println("goroutine ", i, " 看到的计数为：", seeCount, " 当前计数为：", _count, "结束时间为：", time.Now().Format(time.RFC1123Z))
	}

	for i := 1; i <= 30; i++ {
		go reader(i)
	}
	// !!!写锁获取,如果能获取，则执行以下较长时间的代码.
	// !!! 除了作为写入者的主gourouine在运行外，还有30个读取者的gouroutine在并发运行，
	// !!! 主goroutine获得写锁的时候，有多少个读取者已经获得读取锁不确定，有多少个读取者要等到
	// !!! 获得写锁的主goroutine释放写锁后才能得到读取锁的也不确定。
	// !!! 在写锁释放前的这段时间内不会有任何其他线程会读写到count.
	// !!! 也就是说，一旦该所没有释放，其他获取任何锁的线程都会阻塞，等待获得该锁的线程解锁。
	_lock.Lock()
	//!!! 写入者（这里是主gouroutine）获得了写锁之后，必须等待所有的读取锁持者（goroutine）释放读取锁才能开始执行
	//!!! 因此，下一行代码中打印出的时间一定是最后那个读取锁获得者（goroutine）释放了读取锁的时间
	println("在", time.Now().Format(time.RFC1123Z), "之后，main goroutine 将要把count值增加")
	for range 100000000 {
		_count++
	}
	time.Sleep(1 * time.Second)
	println("main goroutine 修改后的count值为：", _count)
	_lock.Unlock() // !!!写锁释放
	_wg.Wait()
	println("----------测试完毕！------------")
}

// TestMakeMapSafeAndGoodReadUsingRWMutex 利用RWMutex的RLock可以多次施加不阻塞特性，
// 提高了一个写操作线程，多个读操作线程的性能。
func TestMakeMapSafeAndGoodReadUsingRWMutex(t *testing.T) {
	var loopWirte = func(m map[int]int, lock *sync.RWMutex) {
		for n := 0; n < 100000; n++ {
			for i := 0; i < 100; i++ {
				lock.Lock()
				m[i] = i
				lock.Unlock()
			}
		}

	}
	var loopRead = func(m map[int]int, lock *sync.RWMutex) {
		for n := 0; n < 100000; n++ {
			lock.RLock()
			for k, v := range m { // 读取map内容
				println(k, "-", v)
			}
			lock.RUnlock()
		}
	}
	lock := &sync.RWMutex{}
	m := make(map[int]int)
	go loopWirte(m, lock)
	go loopRead(m, lock)
	go loopRead(m, lock)
	go loopRead(m, lock)
	//阻塞主线程，除非杀死进程，否则会一直阻塞下去。
	var block = make(chan struct{})
	<-block
}
