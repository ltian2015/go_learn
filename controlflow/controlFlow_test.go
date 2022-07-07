package controlflow

import "testing"

//go中进行debug测试，要求必须安装dlv工具，但dlv工具版本与go版本不匹配时则无法进行测试与debug，
//因此，必须更新dlv工具，但更新安装无法自动覆盖，必须要删除原有已安装的dlv工具，更新方式安装dlv的命令如下：
//sudo rm $GOROOT/bin/dlv
//sudo rm $GOPATH/bin/dlv
// sudo go get -u github.com/go-delve/delve/cmd/dlv

// go的自动化自测试工具go testd对被测试的文件与函数的命名约定：
// 1.测试文件名必须以  _test结尾
// 2.测试函数必须以Test开头，参数必须是(t *testing.T)
// 3.一个包路径下只有一个测试文件。（在vscode工具至少是这样）
//按照上述约定，在vscode中，使用go test -v ${fileDirname}就可以测试当前文件所在包。
//-v 参数可以把测试过程中的fmt.print 结果显示出来。
func TestThePkg(t *testing.T) {
	//TestFor(t)
	//TestSwitch(t)
}
