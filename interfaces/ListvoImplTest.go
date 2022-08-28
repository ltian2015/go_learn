package interfaces

// 里氏替换原则在GO语言中的实现方式与要求。
// 实现抽象类型Reader的一个“具体子类”,该子类Read行为输出了一个比抽象接口要求更具体的类型，
// 应该更可以满足输出的使用者要求。
// 按照里氏替换原则，应该可以实现替换，在Scala语言中，可以该具体实现可以替换抽象概念，但是在GO和Java中都不可以。
// 注意，在GO中，具体类型的方法集必须与抽象类型的方法集完全一致，如果抽象类型的方法集中的方法在
// 签名中使用了抽象类型，那么，具体类型中的方法 也必须使用该抽象类型，而不能使用该抽象类型的”具体实现类型“。
// 这就导致具体类型中必须引用那个抽象类型，使得具体类型依赖来了抽象类型，而不能实现抽象与具体的彻底分离。
// 如果想达成抽象与具体的完全分离，必须使用泛型。用法详见 /generic/seprate_concrete_abstract_test.go

type MyStringReader string

func (ms MyStringReader) Read() string {
	return string(ms)
}
func (ms MyStringReader) Write(content string) {
	ms = MyStringReader(content)
}

// -----------------Writer接口的实现----------------------------
// 该实现意图处理比Writer接口要求更广泛的输入处理能力，从而遵照liskvo替换原则安全替换
// Writer接口，但是确实不允许。而相同的代码在Java中允许。
type MyGenricWriter struct {
	content interface{}
}

func (mgw MyGenricWriter) Write(content interface{}) {
	mgw.content = content
}

func testListov() {
	var ms MyStringReader = "hello"
	// 在Go语言中，认为不可能替换，必须严格与抽象概念要求的类型匹配。而Java语言则可以实现替换。
	//var rd Reader = ms
	println(ms)
	var mgw MyGenricWriter = MyGenricWriter{content: 1000}
	mgw.Write(10000)
	//var wt Writer = mgw //GO
	//wt.Write("1000")
}
