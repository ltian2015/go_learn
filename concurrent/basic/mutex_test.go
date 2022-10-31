package basic_test

import (
	"sync"
	"testing"
)

/*
*
mutex是一个用作互斥锁的结构体(struct)，一个goroutine能够对互斥锁进行锁死（Lock）操作的前提是这个锁还没有锁死，
如果锁已经（被其他goroutine）锁死，必须等待（其他goroutine）开锁操作的完成，
对未锁死的的互斥锁进行开锁操作会抛出panic，而对已锁死的锁进行锁死操作会阻塞。
注意，mutex的零值状态为打开状态。
mutex锁的这种按照开锁、锁死、再开锁、再锁死的这种操作顺序有固定规律的特性，
正好可以用来协调多个线程的执行顺序。

sync.RWMutex ，使得加锁的粒度更细致，读操作与写操作可以分开加锁，
从而可以提高少数线程写，多数线程读的工作场景的性能。
主要是：
它提供了两种加锁方式，”完全锁（Lock）"和"读取锁（RLock）"以及对应的"解完全锁(Unlock)"和”解读取锁(RUnlock)“.

当”完全锁（lock）“操作施加后，再次施加”完全锁(lock)“或"读取锁(Rlock)"操作都会阻塞，
施加加锁操作的线程要等待无任和锁或只有读取锁的状态出现。

当”读取锁（Rlock）“操作施加后，再次施加”读取锁（Rlock）的操作不会阻塞，但施加”完全锁“的操作就会阻塞。
这样，读操作不受影响，同时，又保证了写操作的安全。
*
*/
func TestMutexLock(t *testing.T) {
	var aLock sync.Mutex //Mutex零值处于未锁定的状态，一旦使用就不能再被拷贝。
	var info string

	aLock.Lock()
	go func() {
		info = "hello world"
		aLock.Unlock() //开启锁。
	}()
	aLock.Lock() //如果已锁死，则当前协程必须等待锁的开启。
	println(info)
}

// TryLock试图锁死一个锁，并且返回是否成功。
// 请注意，虽然 TryLock 的正确使用确实存在，但正确使用的情况很少见.
// 并且 TryLock 的使用通常表明在对互斥锁的某些使用中存在更深层次的问题。
// 因为，使用TryLock表明开发者失去了对并行运行中的代码中年加锁，解锁擦洗做顺序的清晰掌控。
// 这是go1.18的加入的新特性，但是几乎无法正确使用，详见：
// https://medium.com/a-journey-with-go/go-story-of-trylock-function-a69ef6dbb410
// https://medium.com/a-journey-with-go/go-mutex-and-starvation-3f4f4e75ad50
func TestMutexTrylock(t *testing.T) {
	var aLock sync.Mutex
	var info string
	aLock.Lock() //
	go func() {
		if aLock.TryLock() { //与else部分完全一样，相当多写了废话，否则不好用。
			info = "hello world"
			aLock.Unlock()
		} else {
			info = "hello world"
			aLock.Unlock()
		}
	}()
	aLock.Lock()
	println(info)
	aLock.Unlock()
}

// 测试map的并发安全性，map并发读写访问是不安全的，可能会抛出如下异常：
// fatal error: concurrent map iteration and map write
// 注意，并发读写冲突发生是随机的，所以只有当并发读写map的次数足够大，或无限循环才能够发生这种冲突。
func TestUnsafeMap(t *testing.T) {
	var loopWirte = func(m map[int]int) {
		for n := 0; n < 100000; n++ { //注意，采用了较大的读写次数
			for i := 0; i < 100; i++ {
				m[i] = i
			}
		}

	}
	var loopRead = func(m map[int]int) {
		for n := 0; n < 100000; n++ { //注意，采用了较大的读写次数
			for k, v := range m {
				println(k, "-", v)
			}
		}

	}
	m := make(map[int]int)
	go loopWirte(m)
	go loopRead(m)
	//阻塞主线程，除非杀死进程，否则会一直阻塞下去。
	var block = make(chan struct{})
	<-block

}

// 使用Mutex使得并发读写map不会发生冲突，
// 注意,mutex加锁使得同时只有一个线程能访问map，无论读写都会排队，而map的读访问是并发安全。
// 这就降低了多线程并发读取map的性能。如果想提高多线程并发读map的性能，需要使用RWMutex
func TestMakeMapSafeUsingMutex(t *testing.T) {
	var loopWirte = func(m map[int]int, lock *sync.Mutex) {
		for n := 0; n < 100000; n++ {
			for i := 0; i < 100; i++ {
				lock.Lock()
				m[i] = i
				lock.Unlock()
			}
		}

	}
	var loopRead = func(m map[int]int, lock *sync.Mutex) {
		for n := 0; n < 100000; n++ {
			lock.Lock()
			for k, v := range m { // 读取map内容
				println(k, "-", v)
			}
			lock.Unlock()
		}
	}
	lock := &sync.Mutex{}
	m := make(map[int]int)
	go loopWirte(m, lock)
	go loopRead(m, lock)
	//阻塞主线程，除非杀死进程，否则会一直阻塞下去。
	var block = make(chan struct{})
	<-block
}

// TestRwmutexLockRlockSequence测试先Lock后Rlock的顺序关系。
// 注意，测试表明，"完全锁"操作(Lock)施加后，无法再施加”读取锁“操作（RLock），线程会阻塞。
func TestRwmutexLockRlockSequence(t *testing.T) {
	var lock = sync.RWMutex{}
	lock.Lock()
	println("locke ok! ")
	lock.RLock() //注意，阻塞在此处，等待Unlock。
	println("Rlock ok! ")
	lock.RUnlock()
	println("RUnlock ok! ")
	lock.Unlock()
	println("Unlock ok! ")
}

// TestRwmutexRlockLockSequence测试先Lock后Rlock的顺序关系。
// 注意，测试表明，"读取锁"操作(RLock)施加后，无法再施加”完全锁“操作（Lock），线程会阻塞。
func TestRwmutexRlockLockSequence(t *testing.T) {
	var lock = sync.RWMutex{}

	lock.RLock()
	println("Rlock ok! ")
	lock.Lock()
	println("locke ok! ") //注意，阻塞在此处，等待RUnlock。
	lock.RUnlock()
	println("RUnlock ok! ")
	lock.Unlock()
	println("Unlock ok! ")
}

// TestRwmutexRlockRLockSequence测试先RLock后，再次RLock的加锁顺序关系。
// 注意，测试表明，"读取锁"操作(RLock)施加后，可以再次施加”读取锁“操作（RLock），线程不会阻塞。
func TestRwmutexRlockRLockSequence(t *testing.T) {
	var lock = sync.RWMutex{}

	lock.RLock()
	println("Rlock ok! ")
	lock.RLock()
	println("Rlock again ok ! ")
	lock.RUnlock()
	println("RUnlock ok! ")
	lock.RUnlock()
	println("RUnlock ok again! ")
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
