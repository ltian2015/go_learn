package singleton

import (
	"sync"
)

/**
  go 语言可以多种方式实现并发安全的单例模式，常见的主要有两种：
  一种是利用包的初始化机制。
  一种是使用sync.Once机制。
**/
//------------下面是采用init机制实现单例----------------------
var config *Config

func init() {
	config = &Config{SystemName: "test system"}
}
func GetConfig() *Config {
	return config
}

type Config struct {
	SystemName string
}

//--------------下面是采用sync.Once机制实现单例------------------
type Config2 struct {
	SystemName string
}

var config2 *Config2

func getConfig() *Config2 {
	var once sync.Once
	//func()是一个闭包函数，在GO中所谓闭包函数是一个捕获了上级函数中可见的（输入、输出）参数变量与局部变量的匿名函数。
	once.Do(func() {
		config2 = &Config2{"test system"}
	})
	return config2
}
