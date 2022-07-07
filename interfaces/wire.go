//+build wireinject

package interfaces

import (
	"github.com/google/wire"
)

//
func InjectStringAppendingFile(id FileID, name FileName, path FilePath, content FileContent) StringAppendingFile {
	//通过“提供者(Providers)” 装配出 “依赖组件” 和 “结果组件”
	wire.Build(NewStringAppendingFile)
	//　伪结果
	return StringAppendingFile{}
}

//
func InjectStringRepalceFile(id FileID, name FileName, path FilePath, content FileContent) StringReplaceFile {
	//通过“提供者(Providers)” 装配出 “依赖组件” 和 “结果组件”
	wire.Build(NewStringReplaceFile)
	//　伪结果
	return StringReplaceFile{}
}
