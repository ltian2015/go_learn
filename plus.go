//用以下自定义构建标记pro(增强的意思)，可以让go build工具决定是否在main包的构建中编译该文件。
//pro是自定义的标记，只要在构建时用go build -tags pro 即可让plus.go文件编译到main包中。
//+build pro

package main

func init() {
	println("增强的特性被加入到了本程序中！")
}



	