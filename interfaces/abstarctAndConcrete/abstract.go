/**接口是GO语言中对行为的抽象，GO之中，只有这一种抽象方式。
   现实世界中，所有事物自身都有属性和行为，
   但是，从GO语言中只有对行为的抽象，也就是说，事物自身的属性必须从外部观察者角度来进行抽象，
   最终在GO语言层面体现为观察者视角的事物属性观察或读取行为的抽象。
   而在传统面向对象语言中，父类可以有自己的属性的归纳和抽象，这与GO语言不同。
   这些传统面向对象语言中，直接对外（观察者）暴露共性属性无法控制读写，外部观察者可以读、也可以写，
   同时，对属性的读写，其实也是两个不同的行为，读是用属性对外部值的赋值，写是外部值对属性的赋值。
   并且通过属性的直接读写无法控制属性读、写行为。如果需要专门控制读写行为，那么往往就会定义供外部观察者访问的对属性的读写行为行数。
   所以，从认识事物的角度看，GO语言更合适，因为抽象是从观察者角度出发的。所谓共性，是观察者观察到的共性，
   而观察者之所以可以观察到共性，那么被观察者一定是以某种展现或操作行为来暴露共性（展现、改变属性的行为或其他行为）。
   抽象是观察者对不同类型的事物的共性的归纳，如果共性属性不能被访问，那么对于观察者来说没有意义，
   提供对共性属性的访问，应该以行为方式暴露，这样有利于隐藏事物内部属性的具体细节。
   人类对事物共性的抽象规律是先看到各种事物，然后观察到这些事物的共性，形成抽象的认识。
   而不是从一无所知开始产生抽象的概念，再看（第一个）给定的事物是否符合抽象。
   所以，从这个角度来看，传统面向对象的设计是不符合人类认识事物的客观规律。
   传统面向对象的类设计的时候，必须先要有抽象的基类，才能定义子类。当我们定义一个类的时候，
   首先要考虑的是，它应该从哪个父类继承，这在对事物认识不全面的时候，很难抉择。
   无疑，GO语言不强制要求对象从某个抽象的父类继承符合人类认识事物共性的规律，很灵活。
   go接口是一组行为的抽象，GO接口支持多组行为通过组合成为更大行为集合。
   同时，实现接口的结构体struct 也可以通过组合来实现更大行为集合的实现。
   另外，GO的这种抽象与实现之间的隐式实现关系对里氏（liskvo）替换原则要求的更严格。

**/

package abstarctandconcrete

//定一个了一个人类可识别的对象的接口规范
type IdentityObject interface {
	GetId() string
	GetName() string
	GetPathName() string
}

//定义了一个可进行字符串读写的对象的接口规范
type IoObject interface {
	Read() interface{}
	Write(interface{}) (int, error)
}

//定义了一个可识别，具有字符串读写能力的对象接口规范
type IdentityStrIoObject interface {
	IdentityObject
	IoObject
}

//定一个了一个文件构建器的规范，这是一个函数类型，通过id,name,path,content构建一个
//由于不支持范型，导致FileBuilder无法让实现方以定义内容类型的方式进行扩展。只能制定为一个可容纳任何数据的通用类型来定义文件的内容类型。
type FileBuilder = func(id, name, path string, content interface{}) IdentityStrIoObject
type FileBuilderStrategy = func(isAppend bool) FileBuilder

//-----------------下面主要研究GO中的里氏（liskvo）替换原则----------------------------
//在GO中，由于没有显式的继承机制，但是有抽象，仍需要考虑具体实现对抽象概念的可替换性，也就是里氏原则。
//GO中，具体实现的行为的输入与输出的类型必须和抽象的行为的输入与输出类型严格匹配才能实现里氏替换原则。
//这与Java，Scala等语言不同。

type Reader interface {
	//该行为要求输出一个通用的类型，意味着任何一种输出类型都可以满足要求。
	//如果GOOGLE支持范型，则这一情况可以获得很大改观。
	Read() interface{}
}
type Writer interface {
	//该行为要求出入一个string类型，意味着本行只能处理string或可以安全替换string类型的“派生类型”。
	//因此，可替换本接口的实现所处理的类型范围必须比string类型更广泛，才能安全替换。
	//在Java中，子类的如果 声明为 Write(Object obj),则可以安全替换本接口。
	//但是，在Go中不可以。
	Write(interface{})
}


