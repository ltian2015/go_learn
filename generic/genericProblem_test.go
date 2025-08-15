package generic

import (
	"fmt"
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

func GetX[P interface{ Point | Rect | Elli }](p P) int {
	return p.X
}

func main() {
	p := Point{1, 2}
	r := Rect{2, 3, 7, 8}
	e := Elli{4, 5, 9, 10}
	fmt.Printf("X: %d %d %d\n", GetX(p), GetX(r), GetX(e))
}
