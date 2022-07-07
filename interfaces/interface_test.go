package interfaces

import (
	"fmt"
	"testing"

	abstarctandconcrete "com.example/golearn/interfaces/abstarctAndConcrete"
)

/**
interfae_test_app这是一个测试应用。该应用通过手工或Wire依赖注入方式将具体实现与抽象进行组装（注入）
因为GO中抽象与实现进行了完全的分离，互不依赖，但是在具体的应用（场景）中，要同时依赖二者。
因为只有抽象，程序无法工作，只有具体实现则无法灵活变更实现的变化。
所以该应用要将具体的实现注入到抽象的概念上，应用的场景逻辑依赖于抽象，不依赖与具体实现。
这样，只要可以注入实现，就可以实现应用（场景）效果的灵活变化。
**/
//需要被注入具体实现的抽象概念的变量。
var identityStrIoObject abstarctandconcrete.IdentityStrIoObject = nil

//
func init() {
	//	manualInject()
	//stringAppendingFileInject()
	stringReplaceFileInject()
}
func manualInject() {
	identityStrIoObject = abstarctandconcrete.NewStringAppendingFile("1000", "lantian.doc", "/neusoft/energy", "lantian is a programer")
}

//通过wire 实现的组装，虽然本程序很简单，体现不出wire自动组装的能力，但是可以说明使用方法。
//wire注入是编译期注入，手写注入的马甲文件（wire.go）,由wire工具自动生成注入实现逻辑wire_gen.go

func stringAppendingFileInject() {
	identityStrIoObject = InjectStringAppendingFile("1000", "lantian.doc", "/neusoft/energy", "lantian is a programer")
}

//通过wire 实现的组装，虽然本程序很简单，体现不出wire自动组装的能力，但是可以说明使用方法。
//wire注入是编译期注入，手写注入的马甲文件（wire.go）,由wire工具自动生成注入实现逻辑wire_gen.go
func stringReplaceFileInject() {
	identityStrIoObject = InjectStringRepalceFile("1000", "lantian.doc", "/neusoft/energy", "lantian is a programer")
}

//测试应用（场景）逻辑，虽然该应用（场景）逻辑很简单，但应用（场景）逻辑中主要对抽象概念进行编程，抽象概念的具体实现则通过依赖注入灵活改变。
func TestInterface(t *testing.T) {
	//不明原因的初始化失败不是调用逻辑的问题，不能返回error，而要抛出异常。
	if identityStrIoObject == nil {
		panic("appliciton initialze failed!")
	}
	fmt.Println(identityStrIoObject.Read())
	identityStrIoObject.Write(" hehe")
	fmt.Println(identityStrIoObject.Read())
}

func TestInterfaceValue(t *testing.T) {
	UseInterfaceValue()
}
func TestEmptyInferfaeValue(t *testing.T) {
	UseEmptyInterface()
}
