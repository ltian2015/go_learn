// go 模块是一组go package的集合。这个“包集合”以go模块作为整体对外发布，
//并通过模块作为访问入口，进行包文件的访问（通过import 指令）。
//一个go.mod文件是一个go 模块的定义文件，go模块与go.mod文件是一一映射关系。
//所以，在go.mod文件中，用来定义模块名称的module 关键字只能出现一次。
//“GO模块”代表了go.mod文件所在路径下的所有go package的集合。
//go module相当于为这些构成这些go package文件的“相对根路径 /”定义了一个唯一标识名，
//知道这个唯一标识名就可以得到所有go package文件的“相对根路径”，
//import 将该“相对根路径/”与具体包文件相对"偏移路径"组合就可以访问到该包的文件。
//这样，其他依赖该包的文件就可以导入包。
//所以，go.mod文件必须放在模块所在的路径下，也就是必须放在“相对根路径/”下。
module com.example/golearn //定义模块的全球唯一标识，一个模块只能有一个go.mod文件，module只能出现一次。

go 1.24 //定义该模块的go语言版本

require (
//golang.org/x/exp/constrains v0.0.0
github.com/google/wire v0.5.0 //定义该模块所依赖的其他模块。一个模块可以依赖很多模块，require可以出现多次。
)
