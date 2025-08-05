package generic

import (
	"fmt"
	"testing"
	"time"
)

//以“类型集合”的形式定义了一个接口。

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Range[P comparable, R any] interface {
	Range(start, end P) R
	DeRange() (start, end P)
	IsIncludedPoint(p P) bool
	IsBeforePoint(p P) bool
	IsAfterPoint(p P) bool
	String() string
	Union(other R) (bool, R)
	UnionOthers(others []R) (bool, R)
	IsIntersected(other R) bool
	Intersect(other R) (bool, R)
	IntersectOthers(other []R) (bool, R)
}
type Interval Range[time.Time, Interval]

func IsIntersected[P comparable, R any](this, other Range[P, R]) bool {
	thisStart, thisEnd := this.DeRange()
	otherStart, otherEnd := other.DeRange()
	isIntersected := (this.IsIncludedPoint(otherStart) || this.IsIncludedPoint(otherEnd) ||
		other.IsIncludedPoint(thisStart) || other.IsIncludedPoint(thisEnd)) &&
		thisEnd != otherStart &&
		otherEnd != thisStart
	return isIntersected
}
func Intersect[P comparable, R any](this, other Range[P, R]) (bool, R) {
	thisStart, thisEnd := this.DeRange()
	otherStart, otherEnd := other.DeRange()
	isIntersected := IsIntersected(this, other)
	var start, end P // 零值
	if !isIntersected {
		return isIntersected, this.Range(start, end)
	}
	start, end = thisStart, thisEnd
	if this.IsIncludedPoint(otherStart) {
		start = otherStart
	}
	if this.IsIncludedPoint(otherEnd) {
		end = otherEnd
	}
	return isIntersected, this.Range(start, end)
}

func Union[P comparable, R any](this, other Range[P, R]) (bool, R) {
	isIntersected := IsIntersected(this, other)
	thisStart, thisEnd := this.DeRange()
	otherStart, otherEnd := other.DeRange()
	isSuccessive := isIntersected || thisStart == otherEnd || thisEnd == otherStart
	start, end := thisStart, thisEnd
	if this.IsAfterPoint(otherStart) {
		start = otherStart
	}
	if this.IsBeforePoint(otherEnd) {
		end = otherEnd
	}
	return isSuccessive, this.Range(start, end)
}

type NumberRange[T number] struct {
	start T
	end   T
}

func CreateNumberRange[T number](p1, p2 T) NumberRange[T] {

	if p1 <= p2 {
		return NumberRange[T]{start: p1, end: p2}
	} else {
		return NumberRange[T]{p2, p1}
	}

}
func (nr NumberRange[T]) Range(start, end T) NumberRange[T] {
	return NumberRange[T]{start: start, end: end}
}

func (nr NumberRange[T]) DeRange() (start, end T) {
	return nr.start, nr.end
}

func (nr NumberRange[T]) IsIncludedPoint(p T) bool {
	return p >= nr.start && p < nr.end
}
func (nr NumberRange[T]) IsBeforePoint(p T) bool {
	return p >= nr.end
}
func (nr NumberRange[T]) IsAfterPoint(p T) bool {
	return p < nr.start
}

func (nr NumberRange[T]) String() string {
	return fmt.Sprintf("NumberRange[%v,%v)", nr.start, nr.end)
}
func (nr NumberRange[T]) IsIntersected(other NumberRange[T]) bool {
	var rgThis Range[T, NumberRange[T]] = nr
	var rgOther Range[T, NumberRange[T]] = other
	return IsIntersected(rgThis, rgOther)
}
func (nr NumberRange[T]) Intersect(other NumberRange[T]) (bool, NumberRange[T]) {
	var rgThis Range[T, NumberRange[T]] = nr
	var rgOther Range[T, NumberRange[T]] = other
	return Intersect(rgThis, rgOther)
}
func (nr NumberRange[T]) IntersectOthers(others []NumberRange[T]) (bool, NumberRange[T]) {
	var result NumberRange[T] = nr
	var isAllIntersected, intersected bool = true, true
	if len(others) == 0 {
		return isAllIntersected, result
	}
	for _, other := range others {
		intersected, result = result.Intersect(other)
		isAllIntersected = isAllIntersected && intersected
	}
	return isAllIntersected, result
}
func (nr NumberRange[T]) Union(other NumberRange[T]) (bool, NumberRange[T]) {
	var rgThis Range[T, NumberRange[T]] = nr
	var rgOther Range[T, NumberRange[T]] = other
	return Union(rgThis, rgOther)
}
func (nr NumberRange[T]) UnionOthers(others []NumberRange[T]) (bool, NumberRange[T]) {
	var result NumberRange[T] = nr
	var isAllSuccessive, successived bool = true, true
	if len(others) == 0 {
		return true, result
	}
	for _, other := range others {
		successived, result = result.Union(other)
		if successived == false {
			isAllSuccessive = false
		}
	}
	return isAllSuccessive, result
}

func TestGenericBasic(t *testing.T) {
	nr1 := CreateNumberRange(1, 6)
	nr2 := CreateNumberRange(3, 5)
	nr3 := CreateNumberRange(3, 10)
	nr4 := CreateNumberRange(11, 15)
	nrs := []NumberRange[int]{nr2, nr3, nr4}
	ok, nr5 := nr1.UnionOthers(nrs)
	var r1, r2, r3, r4 Range[int, NumberRange[int]] = nr1, nr2, nr3, nr4
	_, b := Intersect(nr1, nr2)
	_ = b
	var ok2, nr6 = nr1.IntersectOthers(nrs[:3])
	var isIntersect1, ir1 = Intersect(r1, r2)
	var isIntersect2, ir2 = Intersect(r1, r3)
	var isIntersect3, ir3 = Intersect(r1, r4)
	println("r1: ", r1.String())
	println("r2: ", r2.String())
	println("r3: ", r3.String())
	println("r4: ", r4.String())
	println("r1*r2 : ", ir1.String(), isIntersect1)
	println("r1*r3 : ", ir2.String(), isIntersect2)
	println("r1*r4 : ", ir3.String(), isIntersect3)
	var isSuccessive1, ur1 = Union(r1, r2)
	var isSuccessive2, ur2 = Union(r1, r3)
	var isSuccessive3, ur3 = Union(r1, r4)

	println("r1+r2: ", ur1.String(), isSuccessive1)
	println("r1+r3: ", ur2.String(), isSuccessive2)
	println("r1+r4: ", ur3.String(), isSuccessive3)
	println("r1+..r4 ", nr5.String(), ok)
	println("r1*..r3 ", nr6.String(), ok2)
	for _, v := range nrs[:2] {
		println(v.String())
	}

}

/**


func Except[P comparable](this, other Range[P]) (head, tail Range[P]) {

}
func IsIncluded[P comparable](this, other Range[P]) bool {

}

func IsBefore[P comparable](this, other Range[P]) bool {

}
func IsAfter[P comparable](sp, ep P) bool {

}




//var t time.Time = time.Now()

type NumberRange[E numberic] struct {
	start E
	end   E
}

//泛型类型的方法不能向顶级函数那样定义自己的类型参数列表，如下是无法编译的
/**
func (this Range[T]) Handle3 [E any] (t T, e E) bool {
	return true
}
**/
/**
func (this NumberRange[E]) Include(that NumberRange[E]) bool {
	return this.start <= that.GetStart() && this.end >= that.GetEnd()
}
func (this NumberRange[E]) GetStart() E {
	return this.start
}
func (this NumberRange[E]) GetEnd() E {
	return this.end
}
type TimeInterval struct {
	Start time.Time
	End   time.Time
}

func (this TimeInterval) GetStart() time.Time {
	return this.Start
}
func (this TimeInterval) GetEnd() time.Time {
	return this.End
}

func (this TimeInterval) Include(that TimeInterval) bool {
	return this.Start.Before(that.GetStart()) && this.End.After(that.GetEnd())
}

func (this TimeInterval) Merge(that TimeInterval) TimeInterval {
	var start time.Time
	var end time.Time
	if this.Start.Before(that.GetStart()) {
		start = this.Start
	} else {
		start = that.GetStart()
	}
	if this.End.After(that.GetEnd()) {
		end = this.End
	} else {
		end = that.GetEnd()
	}
	return TimeInterval{Start: start, End: end}
}
**/
