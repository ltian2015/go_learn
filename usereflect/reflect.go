/**
usereflect包主要使用如何使用反射。理解反射的机制和三大法则。
但对于高性能的热点路径代码，最好不要使用反射机制，因为反射目前采用的是运行时库，
内存分配主要在堆中，可能会引起垃圾回收的效率问题。
**/
package usereflect

import (
	"fmt"
	"reflect"
	"time"
	"unsafe"
)

func TestReflect2() {
	//在GO中，每个接口变量都会被编译器翻译为一个类似（VarValue,VarType）这样的一个二元组，
	//其中VarValue被称为动态值，是接口的实际值。VarType被称为动态类型，是变量具体的真实类型，而不是接口类型这样的抽象类型。
	//比如，下面的 reader就是一个（s1,*Student）这样的一个二元组。 *Strudent类型实现了Reader接口的read方法。

	s1 := Student{"lantian"}
	var reader Reader = s1
	fmt.Println(reader.read())
	var x float64 = 3.4
	//GO 反射法则1: 反射机制把接口值（Interface{}）变为反射对象(reflect.Type或reflect.Value)

	fmt.Println("type:", reflect.TypeOf(x).Kind() == reflect.Float64)
	fmt.Println("value:", reflect.ValueOf(x).String())
	var s = Student{"lantian"}
	fmt.Println("type of s : ", reflect.TypeOf(s).Kind())
	//GO 反射法则2， 反射机制可以把反射对象（ reflect.Value）变为接口值(Interface{})
	//再通过接口值（Interface{}）转变为具体类型。
	s2 := reflect.ValueOf(s).Interface().(Student)
	fmt.Println(s2.Name)
	//GO 反射法则3 更改一个反射对象（reflect.Value），那么反射对象值必须可以被设置
	rv := reflect.ValueOf(s1)
	//rv.Set(reflect.ValueOf(Student{"lili"})) //运行时出错
	fmt.Println("rv is settable: ", rv.CanSet())
	rv2 := reflect.ValueOf(&s1) //获得s1的地址拷贝值，赋予rv2.
	fmt.Println("rv2 is settable: ", rv2.CanSet())
	rv3 := rv2.Elem() //rv2是地址拷贝值，其Elem()指向了真实的可改变的元素，也就是s1。
	fmt.Println("rv3 is settable: ", rv3.CanSet())
	rv3.Set(reflect.ValueOf(Student{"lili"}))
	student := rv3.Interface().(Student)
	fmt.Println(student)
	fmt.Println(s1) //此时，s1已经通过反射机制被更改为Student{"lili"}

}

type Reader interface {
	read() string
}
type Student struct {
	Name string
}

func (s Student) read() string {
	return s.Name + " is reading"
}

type Teacher struct {
	Name string
	book string
}

func (this Teacher) read() string {
	return this.Name + " is reading " + this.book
}

func testInterface() {
	s := Student{"lantian"}
	t := Teacher{"lantian", "数学"}
	var reader Reader = s
	fmt.Println(reader.read())
	reader = t
	fmt.Println(reader.read())
	fmt.Println(t.read())
}

//
func LearnRoutine() {
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("i have sleeped 1 second")
	}()
	c := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			c <- i
		}
	}()
	for i := range c {
		fmt.Println(i)
	}
}

//类型转换
func TypeConvert() {
	type IntSlice []int
	type MySlice []int
	var s = []int{}
	var is = IntSlice{}
	var ms = MySlice{}
	var x struct {
		n int `foo:"ok"`
	}
	var y struct {
		n int `bar:"ok"`
	}
	//x = y //x,y两个变量都是未定义的类型，此时需要考虑tag，二者结构体的tag不同，所以不能隐式转换，但是可以显示转换
	//is = ms //两个变量的底层类型相同，但是二者都是已定义类型，所以不能隐式转换。
	//ms = is
	is = IntSlice(ms)
	ms = MySlice(is)
	x = struct {
		n int `foo:"ok"`
	}(y)
	y = struct {
		n int `bar:"ok"`
	}(x)
	s = is //两个变量的底层类型相同，并且 其中有一个变量（s）是未定义类型，则可以隐式转换。
	is = s
	s = ms
	ms = s
	y.n = 5

	typeOfx := reflect.TypeOf(x)

	//可以为if语句声明其作用范围内的变量。好处就是专用于if语句范围内的变量在if语句之外无法访问。
	if fieledType, ok := typeOfx.FieldByName("n"); ok {
		fmt.Println("ok,get field")
		fmt.Println(fieledType.Tag.Get("foo"))
	} else if ok == false {
		fmt.Println("not ok")
	}

}

func add(x, y int) int {
	return x + y
}

func testReflect() {
	type People struct {
		name string
		age  int8
	}
	var p1 = People{"lantian", 47}
	println(p1.name)
	var pt2People = reflect.ValueOf(p1).Pointer()
	var p unsafe.Pointer = unsafe.Pointer(pt2People)
	var p2p *People = (*People)(p)
	p2p.age += 1
	println(p1.age, p2p.age)
}




