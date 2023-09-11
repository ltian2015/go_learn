package basic

import (
	"sync"
	"sync/atomic"
	"testing"
)

//!!! 所谓“原子操作”是指程序中对变量的一个计算操作（比如，加法），会对应 :
//!!! 1. 读变量主存到CPU所在缓存
//!!! 2. 在所在CPU中执行计算
//!!! 3. 回写CPU缓存到主存
//!!! 这三个大的操作步骤，因为CPU架构的不同，细节的操作步骤可能更多，
//!!! 在多CPU计算机中，就可能导致每个CPU对三个操作不同步，使得操作不具备原子性。
//!!! 而“原子操作”可以保证多CPU下变量计算的原子性。
//!!! 原子操作适合“同一个变量，多线程，多写多读”，而java中的“volatile”变量
//!!! 适合“同一变量，多线程，一写多读的场景”。

// 原子的无符号整数类型没有减法操作，只有加法操作，所以，要做加法操作需要一个技巧。
// 这个技巧就是让一个加数的整数减1后按位取反，相当于减去这个加数。
// ^操作符做一元操作符时，是按位取反操作，即，^n 相当于^作为“异或”二元操作符时，m^n操作，
// m与n类型相同，所有位都是1的数。异或二元操作的结果是两个数，如果对应的位的值相异（不同），
// 则结果对应位的值为1，相同为0.
// 一个位上只有1，0，所以有三种情况：1—1；1-0；0—0,对应的按位计算的操作为：
//
//	按位与 &  : 1-1得1，1-0得0，0-0得0
//	按位或 |  : 1-1得1，1-0得1，0-0得0
//	按位异或^ : 1-1得0，1-0得1，0-0得1

func TestUintSub(t *testing.T) {
	var a uint = 100
	var b uint = 10
	c := a + ^(b - 1)
	d := a - b
	e := a + ^b + 1
	println(c, d, e, c == d)
}
func TestNoAutomicAdd(t *testing.T) {
	var number int = 0
	var wg sync.WaitGroup
	//1000个线程同时操作，每个线程对number加1，但由读写number的顺序错乱，导致实际结果小于1000
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			number += 1
			wg.Done()
		}()
	}
	wg.Wait()
	println(number)
}

func TestAutomicAdd(t *testing.T) {

	var number atomic.Int64
	var wg sync.WaitGroup
	//1000个线程同时操作，但是由于采用了原子操作类型，所以最终结果是1000.
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			number.Add(1)
			wg.Done()
		}()
	}
	wg.Wait()
	println(number.Load())
}
