package main

import (
	"errors"
	"log"
	"time"
)

var UserName string
var LoginTime time.Time
var IsLogin bool

//初始化函数，一个文件中可以有多个init()函数。
func init() {
	isLogin, loginTime, err := login("lantian", "123456")
	if err != nil {
		log.Fatal("Login error-" + err.Error())
		return
	}
	IsLogin = isLogin
	LoginTime = loginTime
	UserName = "lantian"
}

//login方法，模拟登录
//函数返回多值时，对返回值进行命名可以提高程序的优雅性。
//命名的返回值变量与输入的参数变量一样，当进入函数时就会被创建，只不过初值都是“零值”。
//由于返回的变量已经存在了，所以return 语句就变得很简单，可以直接return，否则就需要return多个量。
func login(userName string, pwd string) (isLogin bool, loginTime time.Time, err error) {
	if userName == "lantian" && pwd == "123456" {
		isLogin = true
		loginTime = time.Now()
		return
	}
	err = errors.New("bad username or password")
	return
}
