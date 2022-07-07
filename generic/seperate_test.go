package generic

import "testing"

//////////以下代码是站在抽象高度进行编程，完全不知道具体实现的存在，可以独立存在于一个包中/////////////////

/**
https://go.googlesource.com/proposal/+/HEAD/design/43651-type-parameters.md#mutually-referencing-type-parameters
介绍了如何使用泛型对“两族”相互引用的类型进行抽象，从而使得满足这种抽象的具体类型无需知道抽象的存在。
实现了抽象与具体的彻底分离。
而我们也经常会有一些（值）类型的方法（methods）操作了自身，以实现闭合操作。
如果对通过接口来抽象这种操作了自身的类型，比如：
type Mergable interface {
	Merge(other Mergable) Mergable
}
那么，在GO语言中，实现该抽象类型的具体类型的方法的签名必须满足以下形式：
       Merge(other Mergable) Mergable
也就是必须引用抽象类型，而不是自身，比如：
type Score int
func (this Score) Merge(other Mergable) Mergable {
	return ...
}
而我们期望的是具体类型不知道抽象类型的存在，比如：
type Score int
func (this Score) Merge(other Score) Score {
	return ...
}
使用泛型可以表达对这种引用了自身的的类型抽象，而具体类型无需引用抽象的泛型，下面的类型就是具体案例
**/

type Mergable[T any] interface {
	Merge(t T) T
}

/**
   定义对抽象类型进行操作。
**/

func Merge[T any](values []Mergable[T]) T {
	var result T
	for _, v := range values {
		result = v.Merge(result)
	}
	return result
}

//////////以下代码代表两个不同的具体的类型，它们都不知道抽象的存在，可以分别独立存在于一个包中/////////////////
//Score 是一个Mergable类型，它实现了Merge方法。
type Score int

func (this Score) Merge(s Score) Score {
	return Score(this + s)
}

//Power 是一个Mergable类型，它实现了Merge方法。
type Power int

func (this Power) Merge(p Power) Power {
	return Power(this + p)
}

////////////////////以下代码是一个应用，将抽象处理与具体类型相结合////////////////////////
func TestApp(t *testing.T) {
	s := Score(5)
	var ms1 Mergable[Score] = s
	var ms2 Mergable[Score] = Score(6)
	var mss = []Mergable[Score]{ms1, ms2}
	var mr1 = Merge(mss)
	println(mr1)
	var mp1 Mergable[Power] = Power(10)
	var mp2 Mergable[Power] = Power(20)
	var mps = []Mergable[Power]{mp1, mp2}
	mr2 := Merge(mps)
	println(mr2)
}
