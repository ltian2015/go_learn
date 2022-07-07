package usefunction

import (
	"testing"
)

/**
系统内置的函数（位于buildin包与unsafe包中）是一类特殊的函数，它们与程序员自定义的函数有以下不同：
1.系统内置函数不允许向普通自定义函数那样当作函数类型的变量传递。
2.除了copy与recover外，存在返回值的系统内置函数的 返回值不允许被抛弃。
**/
//作为与系统内置函数相对比的一个普通自定义函数。
func square(a int) int {
	return a * a
}
func TestBuildinFunAsVariable(t *testing.T) {
	var _ = square // 普通自定义函数可以被当作变量传递和赋值。
	/** 以下系统内置函数不能作为函数变量传递
		var _ = append
		var _ = copy
		var _ = delete
		var _ = len
		var _ = cap
		var _ = make
		var _ = new
		var _ = complex
		var _ = real
		var _ = imag
		var _ = close
		var _ = panic
		var _ = recover
		var _ = print
		var _ = println
	**/
}
func TestDiscardFuncResult(t *testing.T) {
	square(5)
	srcSlice := []int{1, 2, 3, 4, 5}
	dstSlice := []int{}
	copy(dstSlice, srcSlice)       //内置函数copy允许抛弃返回结果。
	defer copy(dstSlice, srcSlice) //内置函数copy允许抛弃返回结果,可以defer
	recover()                      //内置函数recover允许抛弃返回结果。
	defer recover()                //内置函数recover允许抛弃返回结果,可以defer.

	//new(int) //内置函数new不允许抛弃返回结果。
	//defer new(int) derfer 要求抛弃new函数的返回结果，但new函数不允许抛弃返回结果。
	println("hello")             //无返回结果的内置函数，不存在抛弃返回结果的问题。
	defer println("hello again") //无返回结果的内置函数，不存在抛弃返回结果的问题,可以被defer
}
