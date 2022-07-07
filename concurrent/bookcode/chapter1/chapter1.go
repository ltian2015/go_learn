package chapter1

import (
	"fmt"
	"sync"
	"time"
)

const houseAage int = 100
const data int = houseAage + 100

var BuilderAage int

/**
包初始化，在包被首次加载时被系统自动调用，所以可能会发生在Main函数调用之间。
因为这个函数不应被其它包调用看所以是私有方式声明（首字母小写）
**/
func init() {
	BuilderAage = houseAage + 100
	fmt.Println("dependency packae init")

}

/**
条件竞争就是在并发情况下，多个操作以不确定的顺序操作同一个数据。
导致条件竞争的直接原因是并发执行操作的执行顺序的不确定性所导致的。
而深层次原因则是程序员的顺序化编程思维导致的。
因为，大多数情况下，我们的程序都是按照确定顺序的执行逻辑进行编写。
*/
func RaceCondition() {
	var data int
	go func() {
		data++
		fmt.Printf("the value in sub goroutine is %v \n", data)
	}()

	if data == 0 {
		//由于此时data可能被上面并行的goroutine所更改，所以打印的未必是0，可能是1
		fmt.Printf("the value of data saw by main routine is %v \n", data)
	}
}

//实现了原子操作，原子操作就是在并行执行的情况下，原子操作中的一个多个语句的执行顺序是可以保证的。
//这里有两个临界区，也就有两个原子操作。
//虽然保证了原子操作，但是无法保证两个原子操作的执行的先后顺序。
/**
程序中，“需要或必须”以“独占”方式访问“共享内存”的代码段被称为临界区。
这意味着从进入临界区到离开临界区的代码执行保证原子操作。
（实现原子操作或临界区的常见办法就是加同步锁（ 加锁操作和解锁操作之间的代码就是原子操作，也是），但这也是一种低效的办法，尤其是需要频繁出入临界区的时候，
同时，也要考虑临界区的大小）。

*/
func AtomicOperation() {
	var mutex sync.Mutex
	var data int
	go func() {
		mutex.Lock() //临界区1
		defer mutex.Unlock()
		data++
	}()
	mutex.Lock() //临界区2
	if data == 0 {
		fmt.Printf("the value is %v \n", data)
	} else {
		fmt.Printf("the value is %v \n", data)
	}
	mutex.Unlock()
}

/**
两个（或以上）并发操作(OP1,OP2)都会操作两个（或以上）独立的共享资源(R1，R2），彼此独占有了对方的资源，
又都在等待对方释放独占的资源才能释放继续执行并能已占有资源时，导致谁都无法请求到让自己可以
继续执行下去的资源而互相等待的情况就是死锁。
导致死锁的原因是，这些操作对同一组资源的加锁顺序不同就会导致死锁。
比如 op1的资源加锁顺序是（R1，R2），而Op2的资源加锁顺序是（R2,R1），那么就容易产生死锁。
因为OP1持有了独占了R1，等待独占R2，而同时，OP2 独占了R2，等待独占R1，这就产生了死锁。
保证多个操作以同样的顺序对同样一组资源加锁是避免死锁的根本
*/
func DeadLock() {
	//go支持在函数中声明类型，但不允许直接以func 方式声明子函数
	type shareValue struct {
		mu    sync.Mutex
		value int
	}
	/**
		 等待组类似与Java并发库中的倒数门闩 (CountDownLatch),用于让主例程(或Java线程）
		 等待多个子例程（或线程）的执行完成，然后再统一进行后续的工作。
	     Add(n) 增加等待例程的数量（通常是在主例程中使用）
		 Done() 使得等待数量减1（一般在子例程中使用，当子例程工作完成后，调用该方法使等待数量减1）
		 Wait() 阻塞方法的调用例程（通常是主例程），直到等待数量消减为0.
	*/
	var wg sync.WaitGroup
	/**
	虽然GO不支持像Scala那样，在函数中直接定义或声明子函数，并且子函数可以使用父函数范围内的变量，
	但是下面使用函数字面量并赋值给一个函数该变量，也相当于在函数中定义了子函数
	*/
	printSum := func(v1, v2 *shareValue) {
		defer wg.Done()
		v1.mu.Lock()
		defer v1.mu.Unlock()
		time.Sleep(2 * time.Second)

		v2.mu.Lock()
		defer v2.mu.Unlock()
		fmt.Printf("sum=%v \n", v1.value+v2.value)
	}
	var a, b shareValue
	//
	wg.Add(2)
	go printSum(&a, &b)
	go printSum(&b, &a)
	//调用Wait(）方法的例程将会阻塞，直到等待组的数值消减到0。
	wg.Wait()
}

/**
活锁，
**/
func LiveLock() {
	///cadence:=sync.NewCond(&sync.Mutex{})
}
