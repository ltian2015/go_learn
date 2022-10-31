package basic

import (
	"sync"
	"sync/atomic"
	"testing"
)

// 原子的无符号整数类型没有减法操作，只有加法操作，所以，要做加法操作需要一个技巧。
// 这个技巧就是让一个加数的整数减1后按位取反，相当于减去这个加数。
// ^操作符做一元操作符时，是按位取反操作，即，^n 相当于^作为“异或”二元操作符时，m^n操作，
// m是于类型相同，所有位都是1的数。异或二元操作的结果是两个数，如果对应的位的值相异（不同），
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
	println(c, d, c == d)
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
