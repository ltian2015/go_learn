package normal

/**
本文件主要演示if语句的用法。
  golang中，if语句的语法是：
    if  var declaration;  condition {
    // code to be executed if condition is true
}
与其他语言不同的是，在if语句中可以进行变量声明，其好处是，if 条件表达式 和 if块之外的代码无法访问if句声明的变量。
很多时候，我们需要调用函数之后，立即根据调用结果进行判断（尤其是错误码），并相应的一些处理。
if句声明的变量即能带来方便，又减少了变量的超范围滥用。

**/

import (
	"errors"
	"fmt"
	"testing"
)

// 在go或其他语言中，应尽量减少if else if else的嵌套，
// 尽量在if语句执行体中采用 break, continue, goto, 或 return来减少if else if else 的嵌套
// 这在其他语言中也是通行做法
func TestIf(t *testing.T) {
	//a,s,c无法在if语句之外被访问到
	if a, s, c := 1, "hello world", 2*10; (a+c == 21 || 3 > 0) && (true) {
		println(s)
	}
	//闭包函数，求年龄之和
	sumAges := func(ages ...int) (int, error) {
		var result int = 0
		if len(ages) == 0 {
			return result, nil
		}

		for i, age := range ages {
			if age < 0 {
				errMsg := fmt.Sprintf("bad age %v at the %vnd argument", age, i)
				return 0, errors.New(errMsg)
			}
			result += age
		}
		return result, nil
	}

	if sum, err := sumAges(10, 11, 12, -1); err != nil {
		println(err.Error())
	} else {
		println("all ages together are sum", sum)
	}
}
