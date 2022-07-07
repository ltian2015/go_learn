package interfaces

/**
1.什么是接口类型？
  接口类型定义了一个方法\函数原型（prototype）的集合，换句话说，接口定义了一个方法集，实际上，
  我们可以把接口类型看作一个方法集合，从另外一个角度，也可以将接口类型看作一个行为(behavior)集合。

2.接口的内存结构什么样？

type empty_interface struct {
	dynamicType  *_type         // the dynamic type
	dynamicValue unsafe.Pointer // 被接口值所“装箱(boxing)”的实现了接口其他类型值被称为“动态值（dynamic value）”
}

type non_empty_interface struct {
	dynamicTypeInfo *struct {
		dynamicType *_type       // the dynamic type
		methods     []*_function // method table
	}
	dynamicValue unsafe.Pointer // the dynamic value
}

接口（interface）类型，在Go语言中具有若干重要作用.但是，其首要作用就是接口（interface）类型
使得GO语言支持对值（value）的“装箱（boxing）”。正是这种对值（value）的“装箱（boxing）”功能，
从而使GO语言支持多态和反射。

**/
