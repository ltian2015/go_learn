package goio

import "io"

type MyStringData struct {
	str       string //default ''
	readIndex int    // default 0.
}

const name = "LANTIAN"

var names = []string{"LANTiAN", "LIUFEI", "AWEI"}
var ages = []int{1, 2, 4}

//golang中const只支持基本类型的值（numberic系列，string） ，不支持未命名的字面类型的常量，也不支持自定义的命名类型的产量。
//const names = []string{"LANTiAN", "LIUFEI", "AWEI"}
//const ages = []int{1, 2, 4}
//const ms MyStringData=MyStringData{str:"hello",readIndex:1}
func NewStringData(value string) *MyStringData {
	return &MyStringData{str: value, readIndex: 0}
}
func (msd *MyStringData) Read(p []byte) (n int, err error) {
	//将‘str’转换为字节切片
	var strBytes []byte = []byte(msd.str)
	// 如果已经读取了所有数据，则返回0和io.EOF错误。
	if msd.readIndex >= len(strBytes) {
		return 0, io.EOF
	}
	return 0, io.EOF
}
