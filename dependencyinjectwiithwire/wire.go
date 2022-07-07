//GO build 标记后面必须跟空行，否则GO build工具会将其误认为普通的注释，而非GO build 标记。
//
//+build wireinject
package dependencyinjectwiithwire

import (
	"github.com/google/wire"
)

// 在wire的核心概念中，InitializeEvent是一个注入器（Injector），
//这个注入器把创建Event组件所需要的依赖组件注入进来，然后返回所创建的Event组件。
//当然，wire.go只是真正工作的注入器的伪实现，也就是“注入器马甲（stub）”，
// 它并不会真正的创建Event组件（这里返回一个Event的“零值”,返回任何一个Event值都可以），
//通过wire 命令（工具）编译wire.go文件后会创建“真正注入器”代码文件，也就是wire_gen.go，
//在使用wire命令（工具）之前，wire.go文件必须要通过编译才行。
//在使用wire命令（工具）生成了wire_gen.go之后，由于函数签名相同的“伪注入器”与“真正注入器”同时存在，
//所以必须在构建(编译)的时候使用构建约束（条件编译）//+build wireinject 指明不编译wire.go，而只编译wire_gen.go，
//否则，就会出现重复声明的编译错误，也正因为这样，wire.go文件只会给 wire使用，而不给工程中的包使用，
// 在不删除构建约束//+build wireinject 情况下，不用wire 工具（命令）就不会发现该文件的语法错误。
//在注入器中使用的NewEvent，NewGreeter，NewMessage被成为“提供者（provider）”,
//它们提供所有的组件,包括提供“依赖组件”和被装配出来的“结果组件”。

func InitializeEvent1() Event {
	//通过“提供者(Providers)” 装配出 “依赖组件” 和 “结果组件”
	wire.Build(NewEvent, NewGreeter, NewMessage)
	//　伪结果
	return Event{}
}

//该函数是另一个"注入器(injector)"，该“注入器”所使用的“提供者”函数，
//也就是初始化装配Event所需的各类型的“初始化器（也称提供者，比如，NewEvent2，NewGreeter2）”，
//由于初始化器NewGreeter2不仅需要一个Message类型的参数，还需要一个bool类型的参数，
//Message类型的参数可由NewMessage提供者提供，但是bool类型的参数没有提供者，
//只能由注入器函数自身的函数提供，所以注入器函数多了一个bool类型的isGrumpy参数。

func InitializeEvent2(isGrumpy bool, grumpy bool) (Event, error) {
	//通过“提供者(Providers)” 装配出 “依赖组件” 和 “结果组件”

	wire.Build(NewEvent2, NewGreeter2, NewMessage2)
	//　伪结果
	return Event{}, nil
}





	

