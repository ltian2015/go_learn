package abstarctandconcrete

import (
	"fmt"
)

/**
   这个具体实现文件定义了用字符串文件的符合抽象的文件概念（abstract ）
   注意:
   1.这里行为的接受者类型是结构体的指针类型，而不是结构体类型，
   所以，是结构体的指针类型实现了抽象的概念，而不是结构体类型实现了抽象的概念。
   2. 虽然具体实现符合抽象的所规规范的契约，但实现部分没有显式地依赖抽象。因此，可以独立开发。
      不需要抽象的设计者提供完整的抽象包，才能进行具体实现的开发。
   3. 在应用程序中，会将抽象与具体实现进行组装，也就是依赖注入（手动或用google wire工具）。
**/

//定一个了一个具体的文件描述结构体
type fileProfile struct {
	id   string
	name string
	path string
}

func (fp *fileProfile) GetId() string {
	return fp.id
}

func (fp *fileProfile) GetName() string {
	return fp.name
}

func (fp *fileProfile) GetPathName() string {
	return fp.path + "/" + fp.name
}

//定义一个以字符串追加方式写入string io的结构体
type stringAppendIo struct {
	content string
}

//读取字符串
func (strApdIo *stringAppendIo) Read() interface{} {

	return strApdIo.content
}

//追加方式写入字符串
func (strApdIo *stringAppendIo) Write(content interface{}) (int, error) {
	var err error = nil
	switch strContent := content.(type) {
	case string:
		strApdIo.content += strContent
	default:
		err = fmt.Errorf("not a string")
	}
	return len(strApdIo.content), err
}

type stringReplaceIo struct {
	content string
}

func (strRepIo *stringReplaceIo) Read() interface{} {
	return strRepIo.content
}
func (strRepIo *stringReplaceIo) Write(content interface{}) (int, error) {
	var err error = nil
	switch strContent := content.(type) {
	case string:
		strRepIo.content = strContent
	default:
		err = fmt.Errorf("not a string")
	}
	return len(strRepIo.content), err
}

//将文件描述结构与字符串IO组合成一个字符串内容的文件。
type StringAppendingFile struct {
	//由于被包含的fileProfile、stringAppendIo等接口体通过指针形式的接收器绑定方法，
	//所以外围结构体必须使用指针变量而不是值变量来组合被包含的结构体（fileProfile、stringAppendIo）
	*fileProfile
	*stringAppendIo
	/** 以下书写方式也是可行的，就是使用的时候麻烦
		   fp *fileProfile
		   sai *stringAppendIo
		   这种写法，通过外围包装结构体的实例调用被包装结构体实例内字段或方法的时候需要用被包装结构体的变量名。
		   形如：
		   sf:=stringFile{......}
	       sf.fp.GetId()
		   而采取当前的这种写法则可以简化调用，形如：
		   sf.GetId()
	**/
}
type StringReplaceFile struct {
	*fileProfile
	*stringReplaceIo
}

//由于GO语言的编程思想是抽象与具体实现之间分离，在不同的应用（场景）下将抽象与具体实现进行组装（可使用Wire）。
//因此，GO语言鼓励返回具体的类型，而不是抽象的类型，这样就可以消除具体实现对抽象的引入与依赖。
//有利于实现的调试。
//由于Wire工具的在写注入器的时候，注入器的逻辑是将注入器的输入数据类型与具体提供器的输入类型进行匹配，
//因而，相同类型的多个参数无法匹配到提供器的输入参数，这种情况主要发生在提供器的输入参数类型是
// string ,int 等基本类型的时候，因此，其最佳实践建议，用基本类型作为底层类型，定义个性化的输入参数类型，
//方便注入器进行提供者的输入参数匹配。
type FileID string
type FileName string
type FilePath string
type FileContent string

//这是具体实现的构造函数，从wire角度来看，它们又被成为提供器（Provider）
func NewStringAppendingFile(id FileID, name FileName, path FilePath, content FileContent) StringAppendingFile {
	var fp = fileProfile{id: string(id), name: string(name), path: string(path)}
	var strApndIo = stringAppendIo{content: string(content)}
	return StringAppendingFile{&fp, &strApndIo}
}

//这是具体实现的构造函数，从wire角度来看，它们又被称为提供器（Provider）
func NewStringReplaceFile(id FileID, name FileName, path FilePath, content FileContent) StringReplaceFile {
	var fp = fileProfile{id: string(id), name: string(name), path: string(path)}
	var strRepdIo = stringReplaceIo{content: string(content)}
	return StringReplaceFile{&fp, &strRepdIo}
}
