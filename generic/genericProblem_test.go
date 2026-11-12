package generic

import (
	"fmt"
	"testing"
)

type Point struct {
	X, Y int
}

type Rect struct {
	X, Y, W, H int
}

type Elli struct {
	X, Y, W, H int
}
type PointWithX interface {
	Point | Rect | Elli
}

// !!! 我想看看什么时候这个编译错误可以被编译器解决
func GetX[P PointWithX](p P) int {
	return p.X
}

func main() {
	p := Point{1, 2}
	r := Rect{2, 3, 7, 8}
	e := Elli{4, 5, 9, 10}
	fmt.Printf("X: %d %d %d\n", GetX(p), GetX(r), GetX(e))
}
func Test_Main(t *testing.T) {
	main()
}
