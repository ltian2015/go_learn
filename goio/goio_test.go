package goio

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestBasicIo(t *testing.T) {
	var s = "Hello 中国 is a string"
	r := strings.NewReader(s)
	n, e := io.Copy(os.Stdout, r)
	println("拷贝了", n, "个字节")
	CheckError(e)
	// Hello it is a string
}

func CheckError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
